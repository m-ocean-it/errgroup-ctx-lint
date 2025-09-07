package func_visitor

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type CommentPosition struct {
	Filename string
	Line     int
}

type funcVisitor struct {
	pass        *analysis.Pass
	nolintLines map[CommentPosition]struct{}

	errgroupCtxStack []errgroupCtxStackElement
}

type errgroupCtxStackElement struct {
	o     types.Object
	depth int
}

func New(pass *analysis.Pass, nolintLines map[CommentPosition]struct{}) *funcVisitor {
	return &funcVisitor{
		pass:        pass,
		nolintLines: nolintLines,
	}
}

func (fv *funcVisitor) Visit(node ast.Node, _ bool, stack []ast.Node) bool {
	if node == nil {
		return false
	}

	switch n := node.(type) {
	case *ast.AssignStmt:
		fv.visitAssignStmt(n, len(stack))
	case *ast.DeclStmt:
		fv.visitDeclStmt(n, len(stack))
	case *ast.CallExpr:
		fv.visitCallExpr(n, node)
	}

	return true
}

func (fv *funcVisitor) visitCallExpr(callExpr *ast.CallExpr, node ast.Node) {
	// selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	// if ok && selExpr != nil {
	// 	if selExpr.Sel.Name == "Go" && isPointerToErrgroup(fv.pass.TypesInfo, selExpr.X) {
	// 		fv.maybeCurrentErrGroupGoroutine = node
	// 	}
	// }

	if len(fv.errgroupCtxStack) == 0 {
		return
	}

	for _, arg := range callExpr.Args {
		if !exprIsContext(fv.pass.TypesInfo, arg) {
			continue
		}

		var isInScope bool
		var lastCtx types.Object

		if len(fv.errgroupCtxStack) > 0 {
			lastCtx = fv.errgroupCtxStack[len(fv.errgroupCtxStack)-1].o

			isInScope = true
			// if arg.Pos() > lastCtx.Pos() && arg.Pos() < lastCtx.Pos() {
			// 	isInScope = true
			// }
		}

		if isInScope {
			fIdent, ok := callExpr.Fun.(*ast.SelectorExpr)
			if ok {
				xIdent, ok := fIdent.X.(*ast.Ident)
				if ok {
					if xIdent.Name == "errgroup" || xIdent.Name == "context" {
						return
					}
				}
			}

			argIdent, ok := arg.(*ast.Ident)
			if ok && argIdent != nil {
				obj := fv.pass.TypesInfo.ObjectOf(argIdent)
				if obj != nil {
					if obj.Pos() != lastCtx.Pos() {
						if positionIsNoLint(arg.Pos(), fv.pass.Fset, fv.nolintLines) {
							continue
						}

						fv.pass.Reportf(node.Pos(),
							"passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined")
						// TODO print line of errgroup context
					}
				}
			}
		}
	}
}

func (fv *funcVisitor) visitAssignStmt(assignStmt *ast.AssignStmt, depth int) {
	if len(assignStmt.Rhs) != 1 {
		return
	}

	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok || callExpr == nil {
		return
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr == nil {
		return
	}

	selIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || selIdent == nil {
		return
	}

	if selIdent.Name != "errgroup" {
		return
	}

	// From here, we assume that any context on the left would be an errgroup context.

	for _, leftExpr := range assignStmt.Lhs {
		if exprIsContext(fv.pass.TypesInfo, leftExpr) {
			ctxIdent, ok := leftExpr.(*ast.Ident)
			if !ok || ctxIdent == nil {
				continue
			}

			fv.errgroupCtxStack = append(fv.errgroupCtxStack, errgroupCtxStackElement{
				o:     fv.pass.TypesInfo.ObjectOf(ctxIdent),
				depth: depth,
			})

			break
		}
	}
}

func (fv *funcVisitor) visitDeclStmt(declStmt *ast.DeclStmt, depth int) {
	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	if !ok || genDecl == nil {
		return
	}

	if genDecl.Tok != token.VAR {
		return
	}

LOOP:
	for _, spec := range genDecl.Specs {
		valSpec, ok := spec.(*ast.ValueSpec)
		if !ok || valSpec == nil {
			continue
		}

		if len(valSpec.Values) != 1 {
			continue
		}

		callExpr, ok := valSpec.Values[0].(*ast.CallExpr)
		if !ok || callExpr == nil {
			continue
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selExpr == nil {
			continue
		}

		selIdent, ok := selExpr.X.(*ast.Ident)
		if !ok || selIdent == nil {
			continue
		}

		if selIdent.Name != "errgroup" {
			continue
		}

		// From here, we assume that any context on the left would be an errgroup context.

		for _, leftIdent := range valSpec.Names {
			if exprIsContext(fv.pass.TypesInfo, leftIdent) {
				fv.errgroupCtxStack = append(fv.errgroupCtxStack, errgroupCtxStackElement{
					o:     fv.pass.TypesInfo.ObjectOf(leftIdent),
					depth: depth,
				})

				break LOOP
			}
		}
	}

	return
}

func exprIsContext(typesInfo *types.Info, expr ast.Expr) bool {
	// TODO A more robust approach is probably needed here...

	return typesInfo.TypeOf(expr).String() == "context.Context"
}

func isPointerToErrgroup(typesInfo *types.Info, expr ast.Expr) bool {
	typeOfExpr := typesInfo.TypeOf(expr)

	ptr, ok := typeOfExpr.(*types.Pointer)
	if !ok || ptr == nil {
		return false
	}

	elem := ptr.Elem()

	n, ok := elem.(*types.Named)
	if !ok {
		return false
	}

	o := n.Obj()
	if o == nil {
		return false
	}
	if o.Name() != "Group" {
		return false
	}

	pkg := o.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Name() == "errgroup"
}

func positionIsNoLint(pos token.Pos, fset *token.FileSet, nolintPositions map[CommentPosition]struct{}) bool {
	fullPosition := fset.Position(pos)

	_, isNolint := nolintPositions[CommentPosition{
		Filename: fullPosition.Filename,
		Line:     fullPosition.Line,
	}]

	return isNolint
}
