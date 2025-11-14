package brew

import (
	"errors"
	"fmt"
	"os"

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
			line := fmt.Sprintln(prefix, `"`, item, `"`)
			if _, err := file.Write([]byte(line)); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetupPackages(config brew.Brew, dryRun bool) {
	file, err := os.CreateTemp("", "brew")
	utils.PrintError(err)

	defer file.Close()
	defer os.Remove(file.Name())

	formulaErr := writeBrewFileLines(file, "formula", config.Formulas)
	utils.PrintError(formulaErr)

	caskErr := writeBrewFileLines(file, "cask", config.Casks)
	utils.PrintError(caskErr)

	cmd := fmt.Sprintf("brew bundle --cleanup zap --file=%s", file.Name())

	cmdErr := utils.RunCommand(cmd, "Running homebrew cli", dryRun)
	utils.PrintError(cmdErr)
}
