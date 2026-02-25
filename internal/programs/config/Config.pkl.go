// Code generated from Pkl module `macfigure.config`. DO NOT EDIT.
package config

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/quintisimo/macfigure/internal/programs/brew"
	"github.com/quintisimo/macfigure/internal/programs/cron"
	"github.com/quintisimo/macfigure/internal/programs/dock"
	"github.com/quintisimo/macfigure/internal/programs/home"
	"github.com/quintisimo/macfigure/internal/programs/nsglobaldomain"
	"github.com/quintisimo/macfigure/internal/programs/secret"
)

type Config struct {
	Brew brew.Brew `pkl:"brew"`

	Nsglobaldomain nsglobaldomain.Nsglobaldomain `pkl:"nsglobaldomain"`

	Home []home.Home `pkl:"home"`

	Cron []cron.Cron `pkl:"cron"`

	Dock dock.Dock `pkl:"dock"`

	Secret []secret.Secret `pkl:"secret"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Config
func LoadFromPath(ctx context.Context, path string) (ret Config, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Config
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Config, error) {
	var ret Config
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
