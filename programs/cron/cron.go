package cron

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/utils"
)

type CronProgram struct {
	programs.Program[[]Cron]
}

func (c *CronProgram) Run(logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(c.Input) {
		if removeCronErr := utils.RunCommand("crontab -r", "Removing existing cron jobs", logger, dryRun); removeCronErr != nil {
			return removeCronErr
		}

		cmd := ""
		for _, cron := range c.Input {
			reader, readerErr := utils.ReadFile(cron.Source, logger, dryRun)
			if readerErr != nil {
				return readerErr
			}

			if copyFileErr := utils.WriteFile(reader, cron.Target, logger, dryRun); copyFileErr != nil {
				return copyFileErr
			}

			cmd = fmt.Sprintf("%s\n%s %s", cmd, cron.Schedule, cron.Target)
		}

		cmd = fmt.Sprintf("echo \"%s\" | crontab -", strings.TrimSpace(cmd))
		if cronJobErr := utils.RunCommand(cmd, "Setting up cron jobs", logger, dryRun); cronJobErr != nil {
			return cronJobErr
		}
	}
	return nil
}
