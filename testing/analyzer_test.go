package testing

import (
	"testing"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer"
	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer/func_visitor"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	t.Parallel()

	analysistest.Run(
		t,
		"../testdata/base",
		analyzer.NewAnalyzerWithConfig(func_visitor.Config{
			ErrgroupPackagePaths: []string{
				"github.com/m-ocean-it/errgroup-ctx-lint/testdata/base/errgroup",
			},
		}),
	)
}
