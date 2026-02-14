package secret

import (
	"encoding/base64"
	"fmt"
	"os/exec"
)

// Minimal implementation of keyvault get and set, based on: https://github.com/zalando/go-keyring/blob/master/keyring_darwin.go

const execPathKeychain = "/usr/bin/security"

type KeyvaultItem struct {
	service  string
	username string
}

type KeyvaultOperations interface {
	Get() (string, error)
	Set(password string) error
	Print() error
}

func (k KeyvaultItem) Get() (string, error) {
	encodedPassword, getPasswordErr := exec.Command(
		execPathKeychain,
		"find-generic-password",
		"-s", k.service,
		"-wa", k.username).CombinedOutput()

	if getPasswordErr != nil {
		return "", getPasswordErr
	}

	secret, getPasswordErr := base64.StdEncoding.DecodeString(string(encodedPassword))
	return string(secret), getPasswordErr
}

func (k KeyvaultItem) Set(password string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	if addPasswordErr := exec.Command(
		execPathKeychain,
		"add-generic-password",
		"-U",
		"-s", k.service,
		"-a", k.username,
		"-w", encodedPassword,
	).Run(); addPasswordErr != nil {
		return addPasswordErr
	}

	return nil
}

func (k KeyvaultItem) Print() error {
	password, getPasswordErr := k.Get()
	if getPasswordErr != nil {
		return getPasswordErr
	}

	_, printErr := fmt.Printf("%s: %s\n", k.username, password)
	return printErr
}
