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

	if strings.TrimSpace(*pkgPaths) != "" {
		cfg.ErrgroupPackagePaths = []string{}
		for p := range strings.SplitSeq(*pkgPaths, ",") {
			cfg.ErrgroupPackagePaths = append(
				cfg.ErrgroupPackagePaths,
				strings.TrimSpace(p),
			)
		}
	}

	singlechecker.Main(
		analyzer.NewAnalyzerWithConfig(cfg),
	)
}
