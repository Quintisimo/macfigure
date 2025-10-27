package home

import (
	"fmt"
	"os"

	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/utils"
)

func SetupConfigs(config []home.Home, dryRun bool) {
	for _, item := range config {
		if !dryRun {
			file, err := os.Create(item.Target)
			utils.PrintError(err)

			defer file.Close()
			_, writeErr := file.Write([]byte(item.Content))
			utils.PrintError(writeErr)
		} else {
			utils.DryRunInfo(fmt.Sprintf("Creating file at %s with content:\n%s", item.Target, item.Content))
		}
	}
}
