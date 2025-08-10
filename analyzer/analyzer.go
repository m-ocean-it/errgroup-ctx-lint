package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer/func_visitor"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	nolintDirective = "nolint"
	nolintName      = "errgroup-ctx-lint"
	nolintAll       = "all"
)

var Analyzer = &analysis.Analyzer{
	Name:     "ErrGroupCtxLint",
	Doc:      "Checks that the errgroup's context is passed to operations within the errgroup's goroutines",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	var (
		inspector = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		// nodeFilter  = []ast.Node{(*ast.FuncDecl)(nil)} // TODO optimize?
		nolintLines = getNolintLines(pass.Files, pass.Fset)
	)

	thisFuncVisitor := func_visitor.New(pass, nolintLines)

	inspector.WithStack(nil, thisFuncVisitor.Visit)

	return nil, nil
}

func getNolintLines(files []*ast.File, fset *token.FileSet) map[func_visitor.CommentPosition]struct{} {
	var comments []*ast.CommentGroup
	for _, f := range files {
		comments = append(comments, f.Comments...)
	}

	nolintLines := make(map[func_visitor.CommentPosition]struct{})
	for _, comm := range comments {
		if !commentIsNoLint(comm) {
			continue
		}

		pos := fset.Position(comm.Pos())
		if !pos.IsValid() {
			continue
		}

		nolintLines[func_visitor.CommentPosition{
			Filename: pos.Filename,
			Line:     pos.Line,
		}] = struct{}{}
	}

	return nolintLines
}

func commentIsNoLint(commentGroup *ast.CommentGroup) bool {
	if commentGroup == nil || len(commentGroup.List) == 0 {
		return false
	}

	for _, comm := range commentGroup.List {
		nolintTrimmed := strings.TrimPrefix(comm.Text, "//"+nolintDirective)
		if len(nolintTrimmed) == len(comm.Text) {
			return false
		}

		if nolintTrimmed == "" {
			return true
		}

		colonTrimmed := strings.TrimPrefix(nolintTrimmed, ":")
		if len(colonTrimmed) == len(nolintTrimmed) {
			return false
		}

		nolintList := func() []string {
			list := strings.Split(colonTrimmed, ",")
			for i, linterName := range list {
				list[i] = strings.TrimSpace(linterName)
			}

			return list
		}()

		for _, nolintEntry := range nolintList {
			if nolintEntry == nolintAll || nolintEntry == nolintName {
				return true
			}
		}
	}

	return false
}
