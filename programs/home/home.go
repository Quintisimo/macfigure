package home

import (
	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/gen/home"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/utils"
)

type HomeProgram struct {
	programs.Program[[]home.Home]
}

func (h *HomeProgram) Run(logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(h.Input) {
		for _, item := range h.Input {
			copyFileErr := utils.CopyFile(item.Source, item.Target, logger, dryRun)
			if copyFileErr != nil {
				return copyFileErr
			}
		}
	}
	return nil
}
