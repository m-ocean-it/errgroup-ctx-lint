package func_visitor

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const (
	// TODO respect nolint directive
	nolintDirective = "nolint"
	nolintName      = "errgroup-ctx-lint"
	nolintAll       = "all"
)

type funcVisitor struct {
	pass       *analysis.Pass
	commentMap ast.CommentMap

	maybeErrGroupCtx              types.Object
	maybeCurrentErrGroupGoroutine ast.Node
}

func New(pass *analysis.Pass, commentMap ast.CommentMap) ast.Visitor {
	return &funcVisitor{
		pass:       pass,
		commentMap: commentMap,
	}
}

func (fv *funcVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	return fv.visit(node)
}

func (fv *funcVisitor) visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		return fv.visitCallExpr(n, node)
	case *ast.AssignStmt:
		return fv.visitAssignStmt(n)
	case *ast.DeclStmt:
		return fv.visitDeclStmt(n)
	default:
		return fv
	}
}

func (fv *funcVisitor) visitCallExpr(callExpr *ast.CallExpr, node ast.Node) ast.Visitor {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if ok && selExpr != nil {
		if selExpr.Sel.Name == "Go" && isPointerToErrgroup(fv.pass.TypesInfo, selExpr.X) {
			fv.maybeCurrentErrGroupGoroutine = node
		}
	}

	if fv.maybeErrGroupCtx == nil {
		return fv
	}

	for _, arg := range callExpr.Args {
		if !exprIsContext(fv.pass.TypesInfo, arg) {
			continue
		}

		var isInScope bool
		if fv.maybeCurrentErrGroupGoroutine != nil &&
			arg.Pos() > fv.maybeCurrentErrGroupGoroutine.Pos() &&
			arg.Pos() < fv.maybeCurrentErrGroupGoroutine.End() {

			isInScope = true
		}

		if isInScope {
			argIdent, ok := arg.(*ast.Ident)
			if ok && argIdent != nil {
				obj := fv.pass.TypesInfo.ObjectOf(argIdent)
				if obj != nil {
					if obj.Pos() != fv.maybeErrGroupCtx.Pos() {
						fv.pass.Reportf(node.Pos(),
							"passing non-errgroup context to function withing errgroup-goroutine while there is an errgroup-context defined")
						// TODO print line of errgroup context
					}
				}
			}
		}
	}

	return fv
}

func (fv *funcVisitor) visitAssignStmt(assignStmt *ast.AssignStmt) ast.Visitor {
	if len(assignStmt.Rhs) != 1 {
		return fv
	}

	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok || callExpr == nil {
		return fv
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr == nil {
		return fv
	}

	selIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || selIdent == nil {
		return fv
	}

	if selIdent.Name != "errgroup" {
		return fv
	}

	// From here, we assume that any context on the left would be an errgroup context.

	for _, leftExpr := range assignStmt.Lhs {
		if exprIsContext(fv.pass.TypesInfo, leftExpr) {
			ctxIdent, ok := leftExpr.(*ast.Ident)
			if !ok || ctxIdent == nil {
				continue
			}

			fv.maybeErrGroupCtx = fv.pass.TypesInfo.ObjectOf(ctxIdent)

			break
		}
	}

	return fv
}

func (fv *funcVisitor) visitDeclStmt(declStmt *ast.DeclStmt) ast.Visitor {
	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	if !ok || genDecl == nil {
		return fv
	}

	if genDecl.Tok != token.VAR {
		return fv
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
				fv.maybeErrGroupCtx = fv.pass.TypesInfo.ObjectOf(leftIdent)

				break LOOP
			}
		}
	}

	return fv
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
