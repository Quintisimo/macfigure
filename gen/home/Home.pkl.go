// Code generated from Pkl module `macfigure.home`. DO NOT EDIT.
package home

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Home struct {
	Target string `pkl:"target"`

	Source string `pkl:"source"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Home
func LoadFromPath(ctx context.Context, path string) (ret Home, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Home
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Home, error) {
	var ret Home
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
