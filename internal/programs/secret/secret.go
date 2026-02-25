package secret

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs"
	"github.com/quintisimo/macfigure/internal/utils"
)

type SecretProgram struct {
	programs.Program[[]Secret]
	ExistingSecret *sync.Map
}

func (h *SecretProgram) Run(logger *log.Logger, dryRun bool) error {
	if utils.SliceHasItems(h.Input) {
		for _, item := range h.Input {
			reader, readerErr := DecryptSecretFile(item.Source, logger, dryRun)
			if readerErr != nil {
				return readerErr
			}

			writeFileErr := utils.WriteFile(reader, item.Target, logger, dryRun)
			if writeFileErr != nil {
				return writeFileErr
			}
			h.ExistingSecret.Delete(item.Source)
		}
	}
	return nil
}
