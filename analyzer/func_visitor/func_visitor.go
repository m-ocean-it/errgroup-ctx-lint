package func_visitor

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type CommentPosition struct {
	Filename string
	Line     int
}

type funcVisitor struct {
	pass        *analysis.Pass
	nolintLines map[CommentPosition]struct{}

	errgroupStack errgroupStack
	goStack       goStack
}

type errgroupStack []errgroupStackElement

type errgroupStackElement struct {
	groupObj types.Object
	ctxObj   types.Object
	depth    int
}

func (s errgroupStack) Trim(depth int) errgroupStack {
	if len(s) == 0 {
		return s
	}

	for i, elem := range s {
		if elem.depth > depth {
			return s[:i]
		}
	}

	return s
}

func (s errgroupStack) LastCtx() types.Object {
	for _, elem := range slices.Backward(s) {
		if elem.ctxObj != nil {
			return elem.ctxObj
		}
	}

	return nil
}

type goStack []int

func (gs goStack) Trim(depth int) goStack {
	if len(gs) == 0 {
		return gs
	}

	for i, d := range gs {
		if d > depth {
			return gs[:i]
		}
	}

	return gs
}

func New(pass *analysis.Pass, nolintLines map[CommentPosition]struct{}) *funcVisitor {
	return &funcVisitor{
		pass:          pass,
		nolintLines:   nolintLines,
		errgroupStack: errgroupStack{},
	}
}

func (fv *funcVisitor) Visit(node ast.Node, push bool, stack []ast.Node) bool {
	if node == nil || !push {
		fv.errgroupStack = fv.errgroupStack.Trim(len(stack))
		fv.goStack = fv.goStack.Trim(len(stack))

		return false
	}

	switch n := node.(type) {
	case *ast.AssignStmt:
		fv.visitAssignStmt(n, len(stack))
	case *ast.DeclStmt:
		fv.visitDeclStmt(n, len(stack))
	case *ast.CallExpr:
		fv.visitCallExpr(n, node, len(stack))
	}

	return true
}

func (fv *funcVisitor) visitCallExpr(callExpr *ast.CallExpr, node ast.Node, depth int) {
	if len(fv.errgroupStack) == 0 {
		return
	}

	funSelExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if ok {
		_ = funSelExpr
		xIdent := funSelExpr.X.(*ast.Ident)
		if xIdent != nil {
			xObj := fv.pass.TypesInfo.TypeOf(xIdent)
			xPtr, _ := xObj.(*types.Pointer)
			if xPtr != nil {
				xElem := xPtr.Elem()
				if xElem != nil {
					xNamed, _ := xElem.(*types.Named)
					if xNamed != nil {
						xUnderlyingObj := xNamed.Obj()
						if xUnderlyingObj != nil {
							for _, stackElem := range fv.errgroupStack {
								if stackElem.groupObj.Pos() == xUnderlyingObj.Pos() {
									fv.goStack = append(fv.goStack, depth)

									return
								}
							}
						}
					}
				}
			}
		}
	}

	if len(fv.goStack) == 0 {
		return
	}

	lastCtx := fv.errgroupStack.LastCtx()
	if lastCtx == nil {
		return
	}

	for _, arg := range callExpr.Args {
		if !exprIsContext(fv.pass.TypesInfo, arg) {
			continue
		}

		if len(fv.errgroupStack) > 0 {
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

	selObj := fv.pass.TypesInfo.Uses[selIdent]
	if selObj == nil {
		return
	}
	selPkgName, _ := selObj.(*types.PkgName)
	if selPkgName == nil {
		return
	}
	selPkgNameImported := selPkgName.Imported()
	if selPkgNameImported == nil {
		return
	}

	// TODO: use Path instead of Name?
	if selPkgNameImported.Name() != "errgroup" {
		return
	}

	// From here, we assume that any context on the left would be an errgroup context.

	newErrgroupElement := errgroupStackElement{
		depth: depth,
	}

	var idents []*ast.Ident
	for _, e := range assignStmt.Lhs {
		id, _ := e.(*ast.Ident)
		if id != nil {
			idents = append(idents, id)
		}
	}

	fillStackElemFromIdents(&newErrgroupElement, idents, fv.pass.TypesInfo)

	if newErrgroupElement.groupObj != nil {
		fv.errgroupStack = append(fv.errgroupStack, newErrgroupElement)
	}
}

func fillStackElemFromIdents(elem *errgroupStackElement, idents []*ast.Ident, typesInfo *types.Info) {
	for _, ident := range idents {
		if exprIsContext(typesInfo, ident) {
			elem.ctxObj = typesInfo.ObjectOf(ident)

			break
		}

		leftObj := typesInfo.ObjectOf(ident)
		if leftObj == nil {
			continue
		}

		leftVar, _ := leftObj.(*types.Var)
		if leftVar == nil {
			continue
		}

		leftType := leftVar.Type()
		if leftType == nil {
			continue
		}

		leftPtr, _ := leftType.(*types.Pointer)
		if leftPtr == nil {
			continue
		}

		leftElem := leftPtr.Elem()
		if leftElem == nil {
			continue
		}

		leftNamed, _ := leftElem.(*types.Named)
		if leftNamed == nil {
			continue
		}

		leftUnderlyingObj := leftNamed.Obj()
		if leftUnderlyingObj == nil {
			continue
		}

		if leftUnderlyingObj.Name() == "Group" {
			elem.groupObj = leftUnderlyingObj
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

	newErrgroupElement := errgroupStackElement{
		depth: depth,
	}

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

		fillStackElemFromIdents(&newErrgroupElement, valSpec.Names, fv.pass.TypesInfo)

		if newErrgroupElement.groupObj != nil {
			fv.errgroupStack = append(fv.errgroupStack, newErrgroupElement)

			return
		}
	}
}

func exprIsContext(typesInfo *types.Info, expr ast.Expr) bool {
	// TODO A more robust approach is probably needed here...

	return typesInfo.TypeOf(expr).String() == "context.Context"
}

func positionIsNoLint(pos token.Pos, fset *token.FileSet, nolintPositions map[CommentPosition]struct{}) bool {
	fullPosition := fset.Position(pos)

	_, isNolint := nolintPositions[CommentPosition{
		Filename: fullPosition.Filename,
		Line:     fullPosition.Line,
	}]

	return isNolint
}
