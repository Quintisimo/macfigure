package secret

import (
	"filippo.io/age"
	"github.com/charmbracelet/huh"
)

const serviceName = "macfigure"

var EncryptionKeyItem KeyvaultOperations = &KeyvaultItem{
	service:  serviceName,
	username: "Encryption Key",
}

var DecryptionKeyItem KeyvaultOperations = &KeyvaultItem{
	service:  serviceName,
	username: "Decryption Key",
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

	_, encryptionKeyErr := EncryptionKeyItem.Get()
	_, decryptionKeyErr := DecryptionKeyItem.Get()

	getKeysErr := encryptionKeyErr != nil || decryptionKeyErr != nil

	if getKeysErr {
		regenerateKeys = true
	}

	if !getKeysErr {
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

		decryptionKeySaveErr := DecryptionKeyItem.Set(identity.String())
		if decryptionKeySaveErr != nil {
			return decryptionKeySaveErr
		}

		encryptionKeySaveErr := EncryptionKeyItem.Set(identity.Recipient().String())
		if encryptionKeySaveErr != nil {
			return encryptionKeySaveErr
		}
	}

	return nil
}
