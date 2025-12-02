package envs

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"filippo.io/age"
	"github.com/zalando/go-keyring"
)

const serviceName = "macfigureAge"
const privateKeyName = serviceName + "PrivateKey"
const publicKeyName = serviceName + "PublicKey"

func GetKeys() (publicKey string, privateKey string, error error) {
	publicKey, publicKeyErr := keyring.Get(serviceName, publicKeyName)
	if publicKeyErr != nil {
		return "", "", publicKeyErr
	}

	privateKey, privateKeyErr := keyring.Get(serviceName, privateKeyName)
	if privateKeyErr != nil {
		return "", "", privateKeyErr
	}

	return publicKey, privateKey, nil
}

func GenerateKeys() error {
	_, _, err := GetKeys()
	if err == nil {
		return errors.New("Keys already exist in keychain")
	}

	identity, generationErr := age.GenerateX25519Identity()
	if generationErr != nil {
		return generationErr
	}

	privateKeySaveErr := keyring.Set(serviceName, publicKeyName, identity.Recipient().String())
	if privateKeySaveErr != nil {
		return privateKeySaveErr
	}

	publicKeySaveErr := keyring.Set(serviceName, privateKeyName, identity.String())
	if publicKeySaveErr != nil {
		return publicKeySaveErr
	}

	return nil
}

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
		saveEnvFile("", envFileName, publicKey)
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
