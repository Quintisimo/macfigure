package secret

import (
	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/gen/secret"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/utils"
)

type SecretProgram struct {
	programs.Program[[]secret.Secret]
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
		}
	}
	return nil
}
