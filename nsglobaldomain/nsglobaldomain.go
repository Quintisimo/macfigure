package nsglobaldomain

import (
	"reflect"
	"sync"

	"github.com/quintisimo/macfigure/gen/nsglobaldomain"
	"github.com/quintisimo/macfigure/utils"
)

func WriteConfig(config nsglobaldomain.Nsglobaldomain, dryRun bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	value := reflect.ValueOf(config)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		stringValue := utils.CovertToString(value.Field(i).Elem())

		var cmd string
		if stringValue != "" {
			cmd = "defaults write NSGlobalDomain " + field.Name + " " + stringValue
		} else {
			cmd = "defaults delete -g " + field.Name
		}

		utils.RunCommand(cmd, "Deleting NSGlobalDomain "+field.Name, dryRun)
	}
}
