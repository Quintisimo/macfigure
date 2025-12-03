package cron

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/gen/cron"
	"github.com/quintisimo/macfigure/utils"
)

func SetupCronJobs(crons []cron.Cron, logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(crons) {
		cmd := ""
		for _, cron := range crons {
			copyFileErr := utils.CopyFile(cron.Source, cron.Target, logger, dryRun)
			if copyFileErr != nil {
				return copyFileErr
			}

			cmd = fmt.Sprintf("%s\n%s %s", cmd, cron.Schedule, cron.Target)
		}
		cmd = fmt.Sprintf("echo \"%s\" | crontab -", strings.TrimSpace(cmd))

		cronJobErr := utils.RunCommand(cmd, "Setting up cron jobs", logger, dryRun)
		if cronJobErr != nil {
			return cronJobErr
		}
	}
	return nil
}
