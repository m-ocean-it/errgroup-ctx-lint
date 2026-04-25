package main

import (
	"flag"
	"strings"

	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	pkgPaths := flag.String("pkgs",
		"golang.org/x/sync/errgroup", // Default.
		"Comma-separated list of packages that provide an errgroup. Use in case you're dealing with a non-standard errgroup library.",
	)

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
