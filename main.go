package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/quintisimo/macfigure/brew"
	"github.com/quintisimo/macfigure/cron"
	"github.com/quintisimo/macfigure/dock"
	"github.com/quintisimo/macfigure/gen/config"
	"github.com/quintisimo/macfigure/home"
	"github.com/quintisimo/macfigure/nsglobaldomain"
	"github.com/quintisimo/macfigure/utils"
)

func main() {
	dryRun := flag.Bool("dry-run", true, "Perform a dry run without making any changes")
	configFile := flag.String("config", utils.GetConfigPath(), "Path to the configuration file")
	syncSystem := flag.Bool("sync", false, "Sync system with configuration file")
	flag.Parse()

	if !*syncSystem && !*dryRun {
		flag.PrintDefaults()
		os.Exit(0)
	}

	config, err := config.LoadFromPath(context.Background(), *configFile)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)

	wg.Go(func() {
		brew.SetupPackages(config.Brew, *dryRun)
	})

	wg.Go(func() {
		nsglobaldomain.WriteConfig(config.Nsglobaldomain, *dryRun)
	})

	wg.Go(func() {
		home.SetupConfigs(config.Home, *dryRun)
	})

	wg.Go(func() {
		cron.SetupCronJobs(config.Cron, *dryRun)
	})

	wg.Go(func() {
		dock.SetupDock(config.Dock, *dryRun)
	})

	wg.Wait()
}
