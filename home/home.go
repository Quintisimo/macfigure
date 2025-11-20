package home

import (
	"fmt"
	"os"

	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/utils"
)

func SetupConfigs(config []home.Home, dryRun bool) {
	if utils.SliceHasItems(config) {
		for _, item := range config {
			if !dryRun {
				contents, readErr := os.ReadFile(item.Source)
				utils.PrintError(readErr)

				file, createErr := os.Create(item.Target)
				utils.PrintError(createErr)
				defer file.Close()

				_, writeErr := file.Write(contents)
				utils.PrintError(writeErr)
			} else {
				utils.DryRunInfo(fmt.Sprintf("Creating %s", item.Target))
			}
		}
	}
}
