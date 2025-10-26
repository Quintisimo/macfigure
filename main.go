package main

import (
	"context"
	"flag"

	"github.com/quintisimo/macfigure/brew"
	"github.com/quintisimo/macfigure/gen/config"
	"github.com/quintisimo/macfigure/home"
	"github.com/quintisimo/macfigure/nsglobaldomain"
)

func main() {
	dryRun := *flag.Bool("dry-run", true, "Perform a dry run without making any changes")
	flag.Parse()

	config, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		panic(err)
	}

	brew.SetupPackages(config.Brew, dryRun)
	nsglobaldomain.WriteConfig(config.Nsglobaldomain, dryRun)
	home.SetupConfigs(config.Home, dryRun)
}
