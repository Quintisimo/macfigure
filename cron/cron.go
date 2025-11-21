package cron

import (
	"fmt"
	"strings"

	"github.com/quintisimo/macfigure/gen/cron"
	"github.com/quintisimo/macfigure/utils"
)

func SetupCronJobs(crons []cron.Cron, dryRun bool) {
	if utils.SliceHasItems(crons) {
		cmd := ""
		for _, cron := range crons {
			utils.CopyFile(cron.Source, cron.Target, dryRun)
			cmd = fmt.Sprintf("%s\n%s %s", cmd, cron.Schedule, cron.Target)
		}
		cmd = fmt.Sprintf("echo \"%s\" | crontab -", strings.TrimSpace(cmd))
		utils.RunCommand(cmd, "Setting up cron jobs", dryRun)
	}
}
