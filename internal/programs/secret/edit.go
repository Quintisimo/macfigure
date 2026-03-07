package secret

import (
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/quintisimo/macfigure/internal/utils"
)

func DecryptSecretFile(secretFileName string, logger *log.Logger, dryRun bool) (io.Reader, error) {
	if !dryRun {
		reader, _, decryptErr := decryptFile(secretFileName)
		return reader, decryptErr
	} else {
		utils.DryRunInfo(fmt.Sprintf("Reading %s", secretFileName), logger)
	}
	return nil, nil
}

func decryptFile(secretFileName string) (io.Reader, string, error) {
	decryptionKey, decryptionKeyErr := DecryptionKeyItem.Get()
	if decryptionKeyErr != nil {
		return nil, "", decryptionKeyErr
	}

	identity, identityErr := age.ParseX25519Identity(decryptionKey)
	if identityErr != nil {
		return nil, "", identityErr
	}

	file, fileErr := os.Open(secretFileName)
	if fileErr != nil {
		return nil, "", fileErr
	}
	defer file.Close()

	reader, readerErr := age.Decrypt(file, identity)
	if readerErr != nil {
		return nil, "", readerErr
	}

	content, contentErr := io.ReadAll(reader)
	if contentErr != nil {
		return nil, "", contentErr
	}

	return reader, string(content), nil
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

	_, contents, decryptErr := decryptFile(secretFileName)
	if decryptErr != nil {
		return decryptErr
	}

	if editErr := huh.NewText().Title(secretFileName).Value(&contents).Run(); editErr != nil {
		return editErr
	}

	if saveErr := saveSecretFile(contents, secretFileName); saveErr != nil {
		return saveErr
	}

	return nil
}
