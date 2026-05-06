package analyzer

import (
	"errors"
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
	nolintName      = "errgroupctx"
	nolintAll       = "all"
)

func DefaultConfig() func_visitor.Config {
	return func_visitor.Config{
		ErrgroupPackagePaths: []string{
			"golang.org/x/sync/errgroup",
		},
	}
}

func NewAnalyzerWithConfig(cfg func_visitor.Config) *analysis.Analyzer {
	return newAnalyzer(cfg)
}

func newAnalyzer(cfg func_visitor.Config) *analysis.Analyzer {
	return &analysis.Analyzer{ //nolint:exhaustruct
		Name:     "errgroupctx",
		Doc:      "Checks that errgroup closures use the context derived from a corresponding errgroup",
		Run:      Run(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func Run(cfg func_visitor.Config) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		if !ok {
			return nil, errors.New("unexpectedly type is not *inspector.Inspector")
		}

		var (
			nodeFilter = []ast.Node{
				(*ast.FuncDecl)(nil),
				(*ast.AssignStmt)(nil),
				(*ast.DeclStmt)(nil),
				(*ast.CallExpr)(nil),
			}
			nolintLines = getNolintLines(pass.Files, pass.Fset)
		)

		thisFuncVisitor := func_visitor.New(pass, nolintLines, cfg)

		inspector.WithStack(nodeFilter, thisFuncVisitor.Visit)

		return nil, nil
	}
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
