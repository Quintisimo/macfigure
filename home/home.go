package home

import (
	"fmt"
	"os"
	"sync"

	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/utils"
)

func SetupConfigs(config []home.Home, dryRun bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for _, item := range config {
		if !dryRun {
			file, err := os.Create(item.Target)
			if err != nil {
				panic(err)
			}

			defer file.Close()
			file.Write([]byte(item.Content))
		} else {
			utils.DryRunInfo(fmt.Sprintf("Creating file at %s with content:\n%s", item.Target, item.Content))
		}
	}
}
