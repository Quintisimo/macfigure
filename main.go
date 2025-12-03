package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/lmittmann/tint"
	"github.com/quintisimo/macfigure/brew"
	"github.com/quintisimo/macfigure/cron"
	"github.com/quintisimo/macfigure/dock"
	"github.com/quintisimo/macfigure/gen/config"
	"github.com/quintisimo/macfigure/home"
	"github.com/quintisimo/macfigure/nsglobaldomain"
	"github.com/quintisimo/macfigure/secret"
	"github.com/quintisimo/macfigure/utils"
	"github.com/urfave/cli/v3"
)

func loadConfig(cmd *cli.Command) (config.Config, error) {
	configFile := cmd.String("config")
	return config.LoadFromPath(context.Background(), configFile)
}

func main() {
	defaultLogLevel := "info"
	logLevels := map[string]slog.Level{
		"debug":         slog.LevelDebug,
		defaultLogLevel: slog.LevelInfo,
		"warn":          slog.LevelWarn,
		"error":         slog.LevelError,
	}
	logLevelKeys := slices.Collect(maps.Keys(logLevels))

	cmd := &cli.Command{
		Name:  "macfigure",
		Usage: "A tool to manage macOS configurations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the configuration `file`",
				Value:   utils.GetConfigPath(),
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Value:   true,
				Usage:   "Perform a dry run without making any changes",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "sync",
				Usage:   "Sync system with config",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "loglevel",
						Aliases: []string{"l"},
						Usage:   fmt.Sprintf("Set log level, can be `%s`", strings.Join(logLevelKeys, ", ")),
						Value:   defaultLogLevel,
						Validator: func(s string) error {
							if !slices.Contains(logLevelKeys, s) {
								return fmt.Errorf("invalid log level: %s", s)
							}
							return nil
						},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					dryRun := cmd.Bool("dry-run")

					config, configErr := loadConfig(cmd)
					if configErr != nil {
						return configErr
					}

					logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
						Level: logLevels[cmd.String("loglevel")],
					}))
					createLoggerWithSection := func(section string) *slog.Logger {
						return logger.With(slog.String("section", section))
					}

					var setupErr error
					wg := new(sync.WaitGroup)

					wg.Go(func() {
						brewLogger := createLoggerWithSection("brew")
						setupErr = brew.SetupPackages(config.Brew, brewLogger, dryRun)
					})

					wg.Go(func() {
						nsglobaldomainLogger := createLoggerWithSection("nsglobaldomain")
						setupErr = nsglobaldomain.WriteConfig(config.Nsglobaldomain, nsglobaldomainLogger, dryRun)
					})

					wg.Go(func() {
						homeLogger := createLoggerWithSection("home")
						setupErr = home.SetupConfigs(config.Home, homeLogger, dryRun)
					})

					wg.Go(func() {
						cronLogger := createLoggerWithSection("cron")
						setupErr = cron.SetupCronJobs(config.Cron, cronLogger, dryRun)
					})

					wg.Go(func() {
						dockerLogger := createLoggerWithSection("dock")
						setupErr = dock.SetupDock(config.Dock, dockerLogger, dryRun)
					})

					wg.Wait()
					return setupErr
				},
			},
			{
				Name:  "secret",
				Usage: "Manage secrets with age",
				Commands: []*cli.Command{
					{
						Name:  "generate",
						Usage: "Generate age key and store it in the macOS keychain",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return secret.GenerateKeys()
						},
					},
					{
						Name:  "retrieve",
						Usage: "Retrieve the stored age public key from the macOS keychain",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							publicKey, privateKey, err := secret.GetKeys()
							if err != nil {
								return err
							}

							fmt.Printf("Public key: %s\n", publicKey)
							fmt.Printf("Private key: %s\n", privateKey)
							return nil
						},
					},
					{
						Name:  "edit",
						Usage: "Edit a secrets file",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							config, configErr := loadConfig(cmd)
							if configErr != nil {
								return configErr
							}

							secretPath, secretErr := secret.List(config.Secret)
							if secretErr != nil {
								return secretErr
							}

							if secretPath != "" {
								return secret.Edit(secretPath)
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
