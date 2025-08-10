package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	t.Parallel()

	analysistest.Run(t, "../testdata/base", Analyzer)
}
