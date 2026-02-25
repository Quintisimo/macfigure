package secret

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/utils"
)

func DecryptSecretFile(secretFileName string, logger *log.Logger, dryRun bool) (io.Reader, error) {
	if !dryRun {
		return decryptFile(secretFileName)
	} else {
		utils.DryRunInfo(fmt.Sprintf("Reading %s", secretFileName), logger)
	}
	return nil, nil
}

func decryptFile(secretFileName string) (io.Reader, error) {
	decryptionKey, decryptionKeyErr := DecryptionKeyItem.Get()
	if decryptionKeyErr != nil {
		return nil, decryptionKeyErr
	}

	identity, identityErr := age.ParseX25519Identity(decryptionKey)
	if identityErr != nil {
		return nil, identityErr
	}

	file, fileErr := os.Open(secretFileName)
	if fileErr != nil {
		return nil, fileErr
	}
	defer file.Close()

	reader, readerErr := age.Decrypt(file, identity)
	if readerErr != nil {
		return nil, readerErr
	}
	return reader, nil
}

func openSecretFile(secretFileName string) (string, error) {
	reader, readerErr := decryptFile(secretFileName)
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

func saveSecretFile(contents string, secretFileName string) error {
	encryptionKey, encryptionKeyErr := EncryptionKeyItem.Get()
	if encryptionKeyErr != nil {
		return encryptionKeyErr
	}

	recipient, recipientErr := age.ParseX25519Recipient(encryptionKey)
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
	if _, statErr := os.Stat(secretFileName); os.IsNotExist(statErr) {
		if saveErr := saveSecretFile("", secretFileName); saveErr != nil {
			return saveErr
		}
	}

	tempFileName, tempFileErr := openSecretFile(secretFileName)
	if tempFileErr != nil {
		return tempFileErr
	}
	defer os.Remove(tempFileName)

	contents, editErr := editSecretFile(tempFileName)
	if editErr != nil {
		return editErr
	}

	if saveErr := saveSecretFile(contents, secretFileName); saveErr != nil {
		return saveErr
	}

	return nil
}
