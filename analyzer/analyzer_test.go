package analyzer

import (
	"path/filepath"
	"testing"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer/func_visitor"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Analyzer(t *testing.T) {
	t.Parallel()

	analysistest.Run(
		t,
		filepath.Join(analysistest.TestData(), "base"),
		NewAnalyzerWithConfig(func_visitor.Config{
			ErrgroupPackagePaths: []string{
				"github.com/m-ocean-it/errgroup-ctx-lint/testdata/base/errgroup",
			},
		}),
	)
}
