package programs

import (
	"os"

	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

type Program[Input any] struct {
	Name  string
	Input Input
}

func (p *Program[Input]) GetName() string {
	return p.Name
}

type Execution interface {
	GetName() string
	Run(logger *log.Logger, dryRun bool) error
}

func RunInParallel(programs []Execution, logLevel log.Level, dryRun bool) error {
	logger := log.New(os.Stderr)
	logger.SetLevel(logLevel)

	wg := new(errgroup.Group)

	for _, program := range programs {
		wg.Go(func() error {
			return program.Run(logger.With("section", program.GetName()), dryRun)
		})
	}

	return wg.Wait()
}
