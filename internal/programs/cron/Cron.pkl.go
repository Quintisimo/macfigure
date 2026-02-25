// Code generated from Pkl module `macfigure.cron`. DO NOT EDIT.
package cron

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Cron struct {
	Schedule string `pkl:"schedule"`

	Target string `pkl:"target"`

	Source string `pkl:"source"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Cron
func LoadFromPath(ctx context.Context, path string) (ret Cron, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Cron
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Cron, error) {
	var ret Cron
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
