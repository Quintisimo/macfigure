package env

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"filippo.io/age"
)

func openEnvFile(envFileName string, privateKey string) (string, error) {
	identity, identityErr := age.ParseX25519Identity(privateKey)
	if identityErr != nil {
		return "", identityErr
	}

	file, fileErr := os.Open(envFileName)
	if fileErr != nil {
		return "", fileErr
	}
	defer file.Close()

	reader, readerErr := age.Decrypt(file, identity)
	if readerErr != nil {
		return "", readerErr
	}

	tempFile, tempFileErr := os.CreateTemp("", "env-edit-*.txt")
	if tempFileErr != nil {
		return "", tempFileErr
	}
	defer tempFile.Close()

	if _, copyErr := io.Copy(tempFile, reader); copyErr != nil {
		return "", copyErr
	}
	return tempFile.Name(), nil
}

func editEnvFile(tempFileName string) (string, error) {
	editor := "vim"
	cmd := exec.Command(editor, tempFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if startErr := cmd.Start(); startErr != nil {
		return "", startErr
	}

	if waitErr := cmd.Wait(); waitErr != nil {
		return "", waitErr
	}

	contents, readErr := os.ReadFile(tempFileName)
	if readErr != nil {
		return "", readErr
	}

	return string(contents), nil
}

func saveEnvFile(contents string, envFileName string, publicKey string) error {
	recipient, recipientErr := age.ParseX25519Recipient(publicKey)
	if recipientErr != nil {
		return recipientErr
	}

	env, envErr := os.Create(envFileName)
	if envErr != nil {
		return envErr
	}
	defer env.Close()

	writer, writerErr := age.Encrypt(env, recipient)
	if writerErr != nil {
		return writerErr
	}

	if _, writingErr := io.WriteString(writer, contents); writingErr != nil {
		return writingErr
	}

	if closeErr := writer.Close(); closeErr != nil {
		return closeErr
	}

	return nil
}

func Edit(envFileName string) error {
	publicKey, privateKey, getErr := GetKeys()
	if getErr != nil {
		return getErr
	}

	if _, statErr := os.Stat(envFileName); os.IsNotExist(statErr) {
		fmt.Println("Env file does not exist. Creating new env file:", envFileName)
		if saveErr := saveEnvFile("", envFileName, publicKey); saveErr != nil {
			fmt.Println("Error creating new env file:", envFileName)
			return saveErr
		}
	}

	tempFileName, tempFileErr := openEnvFile(envFileName, privateKey)
	if tempFileErr != nil {
		return tempFileErr
	}
	defer os.Remove(tempFileName)

	contents, editErr := editEnvFile(tempFileName)
	if editErr != nil {
		return editErr
	}

	if saveErr := saveEnvFile(contents, envFileName, publicKey); saveErr != nil {
		return saveErr
	}

	return nil
}
