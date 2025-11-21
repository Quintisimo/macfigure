package nsglobaldomain

import (
	"log/slog"
	"reflect"

	"github.com/quintisimo/macfigure/gen/nsglobaldomain"
	"github.com/quintisimo/macfigure/utils"
)

func WriteConfig(config nsglobaldomain.Nsglobaldomain, logger *slog.Logger, dryRun bool) error {
	writeConfigErr := utils.WriteConfig(reflect.ValueOf(config), "NSGlobalDomain", "defaults write NSGlobalDomain", "defaults delete -g", logger, dryRun)
	if writeConfigErr != nil {
		return writeConfigErr
	}
	return nil
}
