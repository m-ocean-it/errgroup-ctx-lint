package errgroupctxlint

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer"
	"github.com/m-ocean-it/errgroup-ctx-lint/analyzer/func_visitor"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("errgroupctx", New)
}

type Plugin struct {
	settings func_visitor.Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[func_visitor.Config](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: s}, nil
}

func (f *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewAnalyzerWithConfig(f.settings),
	}, nil
}

func (f *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
