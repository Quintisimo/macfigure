package home

import (
	"log/slog"

	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/utils"
)

func SetupConfigs(config []home.Home, logger *slog.Logger, dryRun bool) {
	if utils.SliceHasItems(config) {
		for _, item := range config {
			utils.CopyFile(item.Source, item.Target, logger, dryRun)
		}
	}
}
