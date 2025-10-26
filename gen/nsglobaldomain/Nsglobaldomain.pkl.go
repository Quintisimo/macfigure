// Code generated from Pkl module `macfigure.nsglobaldomain`. DO NOT EDIT.
package nsglobaldomain

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Nsglobaldomain struct {
	AppleIconAppearanceTheme *string `pkl:"AppleIconAppearanceTheme"`

	NSStatusItemSpacing *int `pkl:"NSStatusItemSpacing"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Nsglobaldomain
func LoadFromPath(ctx context.Context, path string) (ret Nsglobaldomain, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Nsglobaldomain
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Nsglobaldomain, error) {
	var ret Nsglobaldomain
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
