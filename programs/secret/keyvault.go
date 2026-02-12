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
	out, err := exec.Command(
		execPathKeychain,
		"find-generic-password",
		"-s", k.service,
		"-wa", k.username).CombinedOutput()

	if err != nil {
		return "", err
	}

	secret, err := base64.StdEncoding.DecodeString(string(out))
	return string(secret), err
}

func (k KeyvaultItem) Set(password string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	cmd := exec.Command(
		execPathKeychain,
		"add-generic-password",
		"-U",
		"-s", k.service,
		"-a", k.username,
		"-w", encodedPassword,
	)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func (k KeyvaultItem) Print() error {
	password, err := k.Get()
	if err != nil {
		return err
	}

	fmt.Printf("%s: %s\n", k.username, password)
	return nil
}
