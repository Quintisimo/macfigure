package nsglobaldomain

import (
	"reflect"

	"github.com/quintisimo/macfigure/gen/nsglobaldomain"
	"github.com/quintisimo/macfigure/utils"
)

func WriteConfig(config nsglobaldomain.Nsglobaldomain, dryRun bool) {
	utils.WriteConfig(reflect.ValueOf(config), "NSGlobalDomain", "defaults write NSGlobalDomain", "defaults delete -g", dryRun)
}
