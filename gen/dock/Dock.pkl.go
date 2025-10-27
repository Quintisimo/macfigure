// Code generated from Pkl module `macfigure.dock`. DO NOT EDIT.
package dock

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Dock struct {
	Apps *[]any `pkl:"apps"`

	Folders *[]string `pkl:"folders"`

	ShowRecents bool `pkl:"showRecents"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Dock
func LoadFromPath(ctx context.Context, path string) (ret Dock, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Dock
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Dock, error) {
	var ret Dock
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
