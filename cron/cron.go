package cron

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/quintisimo/macfigure/gen/cron"
	"github.com/quintisimo/macfigure/utils"
)

func SetupCronJobs(crons []cron.Cron, logger *slog.Logger, dryRun bool) {
	if utils.SliceHasItems(crons) {
		cmd := ""
		for _, cron := range crons {
			utils.CopyFile(cron.Source, cron.Target, logger, dryRun)
			cmd = fmt.Sprintf("%s\n%s %s", cmd, cron.Schedule, cron.Target)
		}
		cmd = fmt.Sprintf("echo \"%s\" | crontab -", strings.TrimSpace(cmd))
		utils.RunCommand(cmd, "Setting up cron jobs", logger, dryRun)
	}
}
