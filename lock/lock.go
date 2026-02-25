package lock

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/quintisimo/macfigure/programs/home"
	"github.com/quintisimo/macfigure/programs/secret"
)

func getPath() (string, error) {
	homeDir := os.Getenv("XDG_CONFIG_HOME")
	if homeDir == "" {
		homeDir = fmt.Sprintf("%s/.config", os.Getenv("HOME"))
	}

	configFolder := fmt.Sprintf("%s/macfigure", homeDir)
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		if mkdirErr := os.Mkdir(configFolder, 0755); mkdirErr != nil {
			return "", mkdirErr
		}
	}
	lockPath := fmt.Sprintf("%s/lock.json", configFolder)

	return lockPath, nil
}

func Create(home []home.Home, secret []secret.Secret) error {
	lockPath, lockPathErr := getPath()
	if lockPathErr != nil {
		return lockPathErr
	}

	lock := make(map[string]string)

	for _, home := range home {
		lock[home.Source] = home.Target
	}

	for _, secret := range secret {
		lock[secret.Source] = secret.Target
	}

	lockFile, lockFileErr := os.Create(fmt.Sprintf("%s/lock.json", lockPath))
	if lockFileErr != nil {
		return lockFileErr
	}
	defer lockFile.Close()

	encoder := json.NewEncoder(lockFile)
	if encodeErr := encoder.Encode(lock); encodeErr != nil {
		return encodeErr
	}

	return nil
}

func Get() (*sync.Map, error) {
	lockPath, lockPathErr := getPath()
	var lock sync.Map

	if lockPathErr != nil {
		return &lock, lockPathErr
	}

	lockFile, lockFileErr := os.Open(fmt.Sprintf("%s/lock.json", lockPath))
	if lockFileErr != nil {
		return &lock, lockFileErr
	}
	defer lockFile.Close()

	decoder := json.NewDecoder(lockFile)
	if decodeErr := decoder.Decode(&lock); decodeErr != nil {
		return &lock, decodeErr
	}

	return &lock, nil
}

func DeleteRemoved(lock *sync.Map) {
	lock.Range(func(_, value any) bool {
		removeErr := os.Remove(value.(string))
		return removeErr == nil
	})
}
