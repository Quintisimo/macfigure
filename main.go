package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	brewConfig "github.com/quintisimo/macfigure/gen/brew"
	"github.com/quintisimo/macfigure/gen/config"
	cronConfig "github.com/quintisimo/macfigure/gen/cron"
	dockConfig "github.com/quintisimo/macfigure/gen/dock"
	homeConfig "github.com/quintisimo/macfigure/gen/home"
	nsglobaldomainConfig "github.com/quintisimo/macfigure/gen/nsglobaldomain"
	secretConfig "github.com/quintisimo/macfigure/gen/secret"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/programs/brew"
	"github.com/quintisimo/macfigure/programs/cron"
	"github.com/quintisimo/macfigure/programs/dock"
	"github.com/quintisimo/macfigure/programs/home"
	"github.com/quintisimo/macfigure/programs/nsglobaldomain"
	"github.com/quintisimo/macfigure/programs/secret"
	"github.com/urfave/cli/v3"
)

func main() {
	defaultLogLevel := "info"
	logLevels := map[string]log.Level{
		"debug":         log.DebugLevel,
		defaultLogLevel: log.InfoLevel,
		"warn":          log.WarnLevel,
		"error":         log.ErrorLevel,
	}
	logLevelKeys := slices.Collect(maps.Keys(logLevels))

	var parsedConfig config.Config

	cli := &cli.Command{
		Name:  "macfigure",
		Usage: "A tool to manage macOS configurations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				Usage:     "Path to the configuration `file`",
				TakesFile: true,
				Required:  true,
				Action: func(ctx context.Context, cmd *cli.Command, path string) error {
					var parsingErr error
					parsedConfig, parsingErr = config.LoadFromPath(context.Background(), path)
					return parsingErr
				},
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
					return spinner.New().
						Title("Applying config...").
						ActionWithErr(func(context.Context) error {
							dryRun := cmd.Bool("dry-run")

							return programs.RunInParallel([]programs.Execution{
								&brew.BrewProgram{
									Program: programs.Program[brewConfig.Brew]{
										Name:  "brew",
										Input: parsedConfig.Brew,
									},
								},
								&nsglobaldomain.NSGlobalDomainProgram{
									Program: programs.Program[nsglobaldomainConfig.Nsglobaldomain]{
										Name:  "nsglobaldomain",
										Input: parsedConfig.Nsglobaldomain,
									},
								},
								&home.HomeProgram{
									Program: programs.Program[[]homeConfig.Home]{
										Name:  "home",
										Input: parsedConfig.Home,
									},
								},
								&cron.CronProgram{
									Program: programs.Program[[]cronConfig.Cron]{
										Name:  "cron",
										Input: parsedConfig.Cron,
									},
								},
								&secret.SecretProgram{
									Program: programs.Program[[]secretConfig.Secret]{
										Name:  "secret",
										Input: parsedConfig.Secret,
									},
								},
								&dock.DockProgram{
									Program: programs.Program[dockConfig.Dock]{
										Name:  "dock",
										Input: parsedConfig.Dock,
									},
								},
							}, logLevels[cmd.String("loglevel")], dryRun)
						}).
						Run()
				},
			},
			{
				Name:  "secret",
				Usage: "Manage secrets with age",
				Commands: []*cli.Command{
					{
						Name:  "keychain",
						Usage: "Manage secrets encryption and decryption key in the macOS keychain",
						Commands: []*cli.Command{
							{
								Name:  "generate",
								Usage: "Generate age key and store it in the macOS keychain",
								Action: func(ctx context.Context, cmd *cli.Command) error {
									return secret.GenerateKeys()
								},
							},
							{
								Name:  "get",
								Usage: "Get the stored age public key from the macOS keychain",
								Action: func(ctx context.Context, cmd *cli.Command) error {
									if encryptionKeyPrintErr := secret.EncryptionKeyItem.Print(); encryptionKeyPrintErr != nil {
										return encryptionKeyPrintErr
									}

									if decryptionKeyPrintErr := secret.DecryptionKeyItem.Print(); decryptionKeyPrintErr != nil {
										return decryptionKeyPrintErr
									}

									return nil
								},
							},
						},
					},
					{
						Name:  "edit",
						Usage: "Edit a secrets file",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							secretPath, secretErr := secret.List(parsedConfig.Secret)
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

	if cliErr := cli.Run(context.Background(), os.Args); cliErr != nil {
		log.Fatal(cliErr)
	}
}
