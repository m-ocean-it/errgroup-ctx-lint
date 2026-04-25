package func_visitor

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type funcVisitor struct {
	cfg Config

	pass        *analysis.Pass
	nolintLines map[CommentPosition]struct{}

	errgroupStack errgroupStack
}

func New(
	pass *analysis.Pass,
	nolintLines map[CommentPosition]struct{},
	cfg Config,
) *funcVisitor {
	if err := cfg.Prepare(); err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	return &funcVisitor{
		cfg:         cfg,
		pass:        pass,
		nolintLines: nolintLines,
	}
}

func (fv *funcVisitor) Visit(node ast.Node, push bool, stack []ast.Node) bool {
	if node == nil || !push {
		fv.errgroupStack = fv.errgroupStack.Trim(len(stack))

		return false
	}

	switch n := node.(type) {
	case *ast.AssignStmt:
		fv.visitAssignStmt(n, len(stack))
	case *ast.DeclStmt:
		fv.visitDeclStmt(n, len(stack))
	case *ast.CallExpr:
		fv.visitCallExpr(n)
	}

	return true
}

func (fv *funcVisitor) visitCallExpr(callExpr *ast.CallExpr) {
	if len(fv.errgroupStack) == 0 {
		return
	}

	errgroupClosure := tryGetErrgroupClosureFromCallExpr(callExpr, fv.pass.TypesInfo, fv.cfg)
	if errgroupClosure == nil {
		return
	}

	sel := callExpr.Fun.(*ast.SelectorExpr) // safe: callIsErrgroupGoOrTryGo verified this
	xIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	recvObj := fv.pass.TypesInfo.ObjectOf(xIdent)
	if recvObj == nil {
		return
	}

	elem := fv.errgroupStack.FindByGroup(recvObj)
	if elem == nil {
		return
	}

	fv.checkClosureForContexts(errgroupClosure, elem)
}

func (fv *funcVisitor) visitAssignStmt(assignStmt *ast.AssignStmt, depth int) {
	if len(assignStmt.Rhs) != 1 {
		return
	}

	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok || callExpr == nil {
		return
	}

	if !callExprPkgIsErrgroup(callExpr, fv.pass.TypesInfo, fv.cfg) {
		return
	}

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

	fillStackElemFromIdents(&newErrgroupElement, idents, fv.pass.TypesInfo, fv.cfg)

	if newErrgroupElement.groupObj != nil {
		fv.errgroupStack = append(fv.errgroupStack, newErrgroupElement)
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

		if !callExprPkgIsErrgroup(callExpr, fv.pass.TypesInfo, fv.cfg) {
			continue
		}

		fillStackElemFromIdents(&newErrgroupElement, valSpec.Names, fv.pass.TypesInfo, fv.cfg)

		if newErrgroupElement.groupObj != nil {
			fv.errgroupStack = append(fv.errgroupStack, newErrgroupElement)

			return
		}
	}
}

func fillStackElemFromIdents(elem *errgroupStackElement, idents []*ast.Ident, typesInfo *types.Info, cfg Config) {
	for _, ident := range idents {
		if ident.Name == "_" {
			continue
		}

		typ := typesInfo.TypeOf(ident)
		if typ != nil && isContextType(typ) {
			elem.ctxObj = typesInfo.ObjectOf(ident)
			elem.ctxName = ident.Name

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

		if leftUnderlyingObj.Name() == "Group" && leftUnderlyingObj.Pkg() != nil && errgroupPkgPathIsEnabled(cfg, leftUnderlyingObj.Pkg().Path()) {
			elem.groupObj = leftObj
		}
	}
}

func isContextType(typ types.Type) bool {
	if typ == nil {
		return false
	}

	named, ok := types.Unalias(typ).(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "context" && obj.Name() == "Context"
}

func positionIsNoLint(pos token.Pos, fset *token.FileSet, nolintPositions map[CommentPosition]struct{}) bool {
	fullPosition := fset.Position(pos)

	_, isNolint := nolintPositions[CommentPosition{
		Filename: fullPosition.Filename,
		Line:     fullPosition.Line,
	}]

	return isNolint
}

func callExprPkgIsErrgroup(callExpr *ast.CallExpr, typesInfo *types.Info, cfg Config) bool {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr == nil {
		return false
	}

	selIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || selIdent == nil {
		return false
	}

	selObj := typesInfo.Uses[selIdent]
	if selObj == nil {
		return false
	}
	selPkgName, _ := selObj.(*types.PkgName)
	if selPkgName == nil {
		return false
	}
	selPkgNameImported := selPkgName.Imported()
	if selPkgNameImported == nil {
		return false
	}

	return errgroupPkgPathIsEnabled(cfg, selPkgNameImported.Path())
}

func errgroupPkgPathIsEnabled(cfg Config, packagePath string) bool {
	return slices.Contains(cfg.ErrgroupPackagePaths, packagePath)
}

func (fv *funcVisitor) checkClosureForContexts(funcLit *ast.FuncLit, elem *errgroupStackElement) {
	if elem.ctxObj == nil {
		return
	}

	closureStart := funcLit.Pos()
	closureEnd := funcLit.End()

	// Identify func lits that are arguments to errgroup Go/TryGo calls,
	// these will be independently analyzed by the inspector, so we skip them
	skipFuncLits := make(map[*ast.FuncLit]struct{})
	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if innerErrgroupClosure := tryGetErrgroupClosureFromCallExpr(call, fv.pass.TypesInfo, fv.cfg); innerErrgroupClosure != nil {
			skipFuncLits[innerErrgroupClosure] = struct{}{}
		}

		return true
	})

	derivedName := elem.ctxName
	if derivedName == "" {
		derivedName = "<errgroup context>"
	}

	// Check all identifiers, skipping nested errgroup callback bodies
	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		if fl, ok := n.(*ast.FuncLit); ok {
			if _, skip := skipFuncLits[fl]; skip {
				return false
			}
		}

		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		obj := fv.pass.TypesInfo.Uses[ident]
		if obj == nil {
			return true
		}

		if _, ok := obj.(*types.Var); !ok {
			return true
		}

		if !isContextType(obj.Type()) {
			return true
		}

		// Allow the errgroup-derived context itself
		if elem.ctxObj != nil && obj == elem.ctxObj {
			return true
		}

		// Allow contexts defined within the closure body
		if obj.Pos() >= closureStart && obj.Pos() < closureEnd {
			return true
		}

		if positionIsNoLint(ident.Pos(), fv.pass.Fset, fv.nolintLines) {
			return true
		}

		fv.pass.Reportf(ident.Pos(),
			"errgroup callback should probably not reference outer context %q, use the errgroup-derived context %q",
			ident.Name, derivedName)

		return true
	})
}

func tryGetErrgroupClosureFromCallExpr(callExpr *ast.CallExpr, typesInfo *types.Info, cfg Config) *ast.FuncLit {
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if sel.Sel.Name != "Go" && sel.Sel.Name != "TryGo" {
		return nil
	}

	obj := typesInfo.Uses[sel.Sel]
	if obj == nil {
		return nil
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return nil
	}

	if fn.Pkg() == nil {
		return nil
	}

	if !errgroupPkgPathIsEnabled(cfg, fn.Pkg().Path()) {
		return nil
	}

	if len(callExpr.Args) != 1 {
		return nil
	}

	funcLit, _ := callExpr.Args[0].(*ast.FuncLit)

	return funcLit
}
