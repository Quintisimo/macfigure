package cron

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/gen/cron"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/utils"
)

type CronProgram struct {
	programs.Program[[]cron.Cron]
}

func (c *CronProgram) Run(logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(c.Input) {
		cmd := ""
		for _, cron := range c.Input {
			reader, readerErr := utils.ReadFile(cron.Source, logger, dryRun)
			if readerErr != nil {
				return readerErr
			}

			copyFileErr := utils.WriteFile(reader, cron.Target, logger, dryRun)
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
