package home

import (
	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/utils"
)

func SetupConfigs(config []home.Home, dryRun bool) {
	if utils.SliceHasItems(config) {
		for _, item := range config {
			utils.CopyFile(item.Source, item.Target, dryRun)
		}
	}
}
