package secret

import (
	"io"
	"os"
	"os/exec"

	"filippo.io/age"
)

func openSecretFile(secretFileName string, privateKey string) (string, error) {
	identity, identityErr := age.ParseX25519Identity(privateKey)
	if identityErr != nil {
		return "", identityErr
	}

	file, fileErr := os.Open(secretFileName)
	if fileErr != nil {
		return "", fileErr
	}
	defer file.Close()

	reader, readerErr := age.Decrypt(file, identity)
	if readerErr != nil {
		return "", readerErr
	}

	tempFile, tempFileErr := os.CreateTemp("", "secret-edit-*.txt")
	if tempFileErr != nil {
		return "", tempFileErr
	}
	defer tempFile.Close()

	if _, copyErr := io.Copy(tempFile, reader); copyErr != nil {
		return "", copyErr
	}
	return tempFile.Name(), nil
}

func editSecretFile(tempFileName string) (string, error) {
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

func saveSecretFile(contents string, secretFileName string, publicKey string) error {
	recipient, recipientErr := age.ParseX25519Recipient(publicKey)
	if recipientErr != nil {
		return recipientErr
	}

	secret, secretErr := os.Create(secretFileName)
	if secretErr != nil {
		return secretErr
	}
	defer secret.Close()

	writer, writerErr := age.Encrypt(secret, recipient)
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

func Edit(secretFileName string) error {
	publicKey, privateKey, getErr := GetKeys()
	if getErr != nil {
		return getErr
	}

	if _, statErr := os.Stat(secretFileName); os.IsNotExist(statErr) {
		if saveErr := saveSecretFile("", secretFileName, publicKey); saveErr != nil {
			return saveErr
		}
	}

	tempFileName, tempFileErr := openSecretFile(secretFileName, privateKey)
	if tempFileErr != nil {
		return tempFileErr
	}
	defer os.Remove(tempFileName)

	contents, editErr := editSecretFile(tempFileName)
	if editErr != nil {
		return editErr
	}

	if saveErr := saveSecretFile(contents, secretFileName, publicKey); saveErr != nil {
		return saveErr
	}

	return nil
}
