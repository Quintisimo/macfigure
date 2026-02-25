package home

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs"
	"github.com/quintisimo/macfigure/internal/utils"
)

type HomeProgram struct {
	programs.Program[[]Home]
	ExisitingHome *sync.Map
}

func (h *HomeProgram) Run(logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(h.Input) {
		for _, item := range h.Input {
			reader, readerErr := utils.ReadFile(item.Source, logger, dryRun)
			if readerErr != nil {
				return readerErr
			}

			if copyFileErr := utils.WriteFile(reader, item.Target, logger, dryRun); copyFileErr != nil {
				return copyFileErr
			}
			h.ExisitingHome.Delete(item.Source)
		}
	}
	return nil
}
