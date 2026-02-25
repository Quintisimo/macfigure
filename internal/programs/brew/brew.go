package brew

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs"
	"github.com/quintisimo/macfigure/internal/utils"
)

type BrewProgram struct {
	programs.Program[Brew]
}

func writeBrewFileLines(file *os.File, packagesType string, packages []string) error {
	prefix := "brew"

	if packagesType != "cask" && packagesType != "formula" {
		return errors.New("Invalid type: must be 'cask' or 'formula'")
	}

	if packagesType == "cask" {
		prefix = "cask"
	}

	if utils.SliceHasItems(packages) {
		for _, item := range packages {
			line := fmt.Sprintln(prefix, `"`, item, `"`)
			if _, writeErr := file.Write([]byte(line)); writeErr != nil {
				return writeErr
			}
		}
	}
	return nil
}

func (b *BrewProgram) Run(logger *log.Logger, dryRun bool) error {
	file, fileErr := os.CreateTemp("", "brewfile-*.Brewfile")
	if fileErr != nil {
		return fileErr
	}

	defer file.Close()
	defer os.Remove(file.Name())

	if formulaErr := writeBrewFileLines(file, "formula", b.Input.Formulas); formulaErr != nil {
		return formulaErr
	}

	if caskErr := writeBrewFileLines(file, "cask", b.Input.Casks); caskErr != nil {
		return caskErr
	}

	cmd := fmt.Sprintf("brew bundle --cleanup zap --file=%s", file.Name())
	if cmdErr := utils.RunCommand(cmd, "Running homebrew cli", logger, dryRun); cmdErr != nil {
		return cmdErr
	}
	return nil
}
