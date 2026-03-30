package lock

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs/home"
	"github.com/quintisimo/macfigure/internal/programs/secret"
	"github.com/quintisimo/macfigure/internal/utils"
)

func createLockLogger(logLevel log.Level) *log.Logger {
	return utils.CreateLogger(logLevel, "lock")
}

func getLockPath(configDir string, logLevel log.Level, dryRun bool) (string, error) {
	lockPath := fmt.Sprintf("%s/config-lock.json", configDir)

	if dryRun {
		utils.DryRunInfo(fmt.Sprintf("Lock file path: %s", lockPath), createLockLogger(logLevel))
	}

	return lockPath, nil
}

func Create(home []home.Home, secret []secret.Secret, configDir string, logLevel log.Level, dryRun bool) error {
	lockPath, lockPathErr := getLockPath(configDir, logLevel, dryRun)
	if lockPathErr != nil {
		return lockPathErr
	}

	if !dryRun {
		lock := make(map[string]string)

		for _, home := range home {
			lock[home.Source] = home.Target
		}

		for _, secret := range secret {
			lock[secret.Source] = secret.Target
		}

		lockFile, lockFileErr := os.Create(lockPath)
		if lockFileErr != nil {
			return lockFileErr
		}
		defer lockFile.Close()

		encoder := json.NewEncoder(lockFile)
		if encodeErr := encoder.Encode(lock); encodeErr != nil {
			return encodeErr
		}
	} else {
		utils.DryRunInfo("Write lock file with current state", createLockLogger(logLevel))
	}

	return nil
}

func Get(configDir string, logLevel log.Level, dryRun bool) (*sync.Map, error) {
	lockPath, lockPathErr := getLockPath(configDir, logLevel, dryRun)
	var lock sync.Map

	if !dryRun {
		if lockPathErr != nil {
			return &lock, lockPathErr
		}

		if _, statErr := os.Stat(lockPath); os.IsNotExist(statErr) {
			return &lock, nil
		} else if statErr != nil {
			return &lock, statErr
		}

		lockFile, lockFileErr := os.Open(lockPath)
		if lockFileErr != nil {
			return &lock, lockFileErr
		}
		defer lockFile.Close()

		decoder := json.NewDecoder(lockFile)
		if decodeErr := decoder.Decode(&lock); decodeErr != nil {
			return &lock, decodeErr
		}
	} else {
		utils.DryRunInfo("Read lock file if exists", createLockLogger(logLevel))
	}

	return &lock, nil
}

func DeleteRemoved(lock *sync.Map, logLevel log.Level, dryRun bool) {
	if !dryRun {
		lock.Range(func(_, value any) bool {
			removeErr := os.Remove(value.(string))
			return removeErr == nil
		})
	} else {
		utils.DryRunInfo("Delete removed files", createLockLogger(logLevel))
	}
}
