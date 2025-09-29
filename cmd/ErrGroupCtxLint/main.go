package main

import (
	_ "flag"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(
		analyzer.NewAnalyzerWithConfig(
			analyzer.DefaultConfig, // TODO: overwrite from command-line args
		),
	)
}
