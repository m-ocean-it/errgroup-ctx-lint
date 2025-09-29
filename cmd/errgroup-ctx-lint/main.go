package main

import (
	"flag"
	"strings"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	pkgPaths := flag.String("pkg_paths", "", "TODO") // TODO

	flag.Parse()

	cfg := analyzer.DefaultConfig

	if *pkgPaths != "" {
		cfg.ErrgroupPackagePaths = strings.Split(*pkgPaths, ",")
	}

	singlechecker.Main(
		analyzer.NewAnalyzerWithConfig(cfg),
	)
}
