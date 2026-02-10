package nsglobaldomain

import (
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/gen/nsglobaldomain"
	"github.com/quintisimo/macfigure/programs"
	"github.com/quintisimo/macfigure/utils"
)

type NSGlobalDomainProgram struct {
	programs.Program[nsglobaldomain.Nsglobaldomain]
}

func (n *NSGlobalDomainProgram) Run(logger *log.Logger, dryRun bool) error {
	writeConfigErr := utils.WriteConfig(reflect.ValueOf(n.Input), "NSGlobalDomain", "defaults write NSGlobalDomain", "defaults delete -g", logger, dryRun)
	if writeConfigErr != nil {
		return writeConfigErr
	}
	return nil
}
