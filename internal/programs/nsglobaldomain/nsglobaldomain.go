package nsglobaldomain

import (
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs"
	"github.com/quintisimo/macfigure/internal/utils"
)

type NSGlobalDomainProgram struct {
	programs.Program[Nsglobaldomain]
}

func (n *NSGlobalDomainProgram) Run(logger *log.Logger, dryRun bool) error {
	if writeConfigErr := utils.WriteConfig(reflect.ValueOf(n.Input), "NSGlobalDomain", "defaults write NSGlobalDomain", "defaults delete -g", logger, dryRun); writeConfigErr != nil {
		return writeConfigErr
	}
	return nil
}
