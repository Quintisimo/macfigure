// Code generated from Pkl module `macfigure.brew`. DO NOT EDIT.
package brew

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Brew struct {
	Casks []string `pkl:"casks"`

	Formulas []string `pkl:"formulas"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Brew
func LoadFromPath(ctx context.Context, path string) (ret Brew, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Brew
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Brew, error) {
	var ret Brew
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
