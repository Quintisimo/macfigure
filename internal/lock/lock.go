package lock

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/programs/home"
	"github.com/quintisimo/macfigure/internal/programs/secret"
	"github.com/quintisimo/macfigure/internal/utils"
)

func createLockLogger(logLevel log.Level) *log.Logger {
	return utils.CreateLogger(logLevel, "lock")
}

func getLockPath(logLevel log.Level, dryRun bool) (string, error) {
	homeDir := os.Getenv("XDG_CONFIG_HOME")
	if homeDir == "" {
		homeDir = fmt.Sprintf("%s/.config", os.Getenv("HOME"))
	}

	lockPath := fmt.Sprintf("%s/macfigure/lock.json", homeDir)

	if !dryRun {
		if mkDirAllErr := os.MkdirAll(filepath.Dir(lockPath), 0755); mkDirAllErr != nil {
			return "", mkDirAllErr
		}
	} else {
		utils.DryRunInfo("Create lock file if not exists", createLockLogger(logLevel))
	}

	return lockPath, nil
}

func Create(home []home.Home, secret []secret.Secret, logLevel log.Level, dryRun bool) error {
	lockPath, lockPathErr := getLockPath(logLevel, dryRun)
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

func Get(logLevel log.Level, dryRun bool) (*sync.Map, error) {
	lockPath, lockPathErr := getLockPath(logLevel, dryRun)
	var lock sync.Map

	if !dryRun {
		if lockPathErr != nil {
			return &lock, lockPathErr
		}

		if _, statErr := os.Stat(lockPath); os.IsNotExist(statErr) {
			return &lock, nil
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
