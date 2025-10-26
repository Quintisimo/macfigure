package brew

import (
	"errors"
	"os"
	"sync"

	"github.com/quintisimo/macfigure/gen/brew"
	"github.com/quintisimo/macfigure/utils"
)

func writeBrewFileLines(file *os.File, packagesType string, packages *[]string) error {
	prefix := "brew"

	if packagesType != "cask" && packagesType != "formula" {
		return errors.New("Invalid type: must be 'cask' or 'formula'")
	}

	if packagesType == "cask" {
		prefix = "cask"
	}

	if packages != nil {
		for _, item := range *packages {
			line := prefix + " " + "\"" + item + "\"" + "\n"
			if _, err := file.Write([]byte(line)); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetupPackages(config brew.Brew, dryRun bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	file, err := os.CreateTemp("", "brew")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	writeBrewFileLines(file, "formula", config.Formulas)
	writeBrewFileLines(file, "cask", config.Casks)

	cmd := "brew bundle --cleanup zap --file=" + file.Name()

	utils.RunCommand(cmd, "Running homebrew cli", dryRun)
}
