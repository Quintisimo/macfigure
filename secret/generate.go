package secret

import (
	"filippo.io/age"
	"github.com/charmbracelet/huh"
)

const serviceName = "macfigure"

var encryptionKeyItem KeyvaultOperations = &KeyvaultItem{
	service:  serviceName,
	username: "encryptionKey",
}

var decryptionKeyItem KeyvaultOperations = &KeyvaultItem{
	service:  serviceName,
	username: "decryptionKey",
}

func GetKeys() (encryptionKey string, decryptionKey string, error error) {
	encryptionKey, encryptionKeyErr := encryptionKeyItem.Get()
	if encryptionKeyErr != nil {
		return "", "", encryptionKeyErr
	}

	decryptionKey, decryptionKeyErr := decryptionKeyItem.Get()
	if decryptionKeyErr != nil {
		return "", "", decryptionKeyErr
	}

	return encryptionKey, decryptionKey, nil
}

func confirmDialog(title string, value *bool) error {
	confirmErr := huh.NewConfirm().
		Title(title).
		Affirmative("Yes!").
		Negative("No.").
		Value(value).
		Run()

	if confirmErr != nil {
		return confirmErr
	}

	return nil
}

func GenerateKeys() error {
	regenerateKeys := false

	_, _, getKeysErr := GetKeys()

	if getKeysErr != nil {
		regenerateKeys = true
	}

	if getKeysErr == nil {
		confirmErr := confirmDialog("Secret keys already exist in the keychain, do you want to regenerate them?", &regenerateKeys)

		if confirmErr != nil {
			return confirmErr
		}

		if regenerateKeys {
			confirmAgainErr := confirmDialog("Are you sure? If you currently have secrets encrypted with the old keys you will not be able to decrypt them", &regenerateKeys)

			if confirmAgainErr != nil {
				return confirmAgainErr
			}
		}
	}

	if regenerateKeys {
		identity, generationErr := age.GenerateX25519Identity()
		if generationErr != nil {
			return generationErr
		}

		decryptionKeySaveErr := decryptionKeyItem.Set(identity.String())
		if decryptionKeySaveErr != nil {
			return decryptionKeySaveErr
		}

		encryptionKeySaveErr := encryptionKeyItem.Set(identity.Recipient().String())
		if encryptionKeySaveErr != nil {
			return encryptionKeySaveErr
		}
	}

	return nil
}
