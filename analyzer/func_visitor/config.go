package func_visitor

import "errors"

const DefaultPkgPath = "golang.org/x/sync/errgroup"

type Config struct {
	ErrgroupPackagePaths []string
}

func (c *Config) Prepare() error {
	if c == nil {
		return errors.New("config is nil")
	}

	if len(c.ErrgroupPackagePaths) == 0 {
		// TODO: log

		c.ErrgroupPackagePaths = []string{
			DefaultPkgPath,
		}
	}

	return nil
}
