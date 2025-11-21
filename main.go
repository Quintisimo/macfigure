package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/quintisimo/macfigure/brew"
	"github.com/quintisimo/macfigure/cron"
	"github.com/quintisimo/macfigure/dock"
	"github.com/quintisimo/macfigure/gen/config"
	"github.com/quintisimo/macfigure/home"
	"github.com/quintisimo/macfigure/nsglobaldomain"
	"github.com/quintisimo/macfigure/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "macfigure",
		Usage: "A tool to manage macOS configurations",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: true,
				Usage: "Perform a dry run without making any changes",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "sync",
				Usage: "Sync system with config",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config",
						Usage: "Path to the configuration file",
						Value: utils.GetConfigPath(),
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					dryRun := cmd.Bool("dry-run")
					configFile := cmd.String("config")

					config, err := config.LoadFromPath(context.Background(), configFile)
					if err != nil {
						panic(err)
					}

					wg := new(sync.WaitGroup)

					wg.Go(func() {
						brew.SetupPackages(config.Brew, dryRun)
					})

					wg.Go(func() {
						nsglobaldomain.WriteConfig(config.Nsglobaldomain, dryRun)
					})

					wg.Go(func() {
						home.SetupConfigs(config.Home, dryRun)
					})

					wg.Go(func() {
						cron.SetupCronJobs(config.Cron, dryRun)
					})

					wg.Go(func() {
						dock.SetupDock(config.Dock, dryRun)
					})

					wg.Wait()
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
