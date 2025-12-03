package secret

import (
	"filippo.io/age"
	"github.com/charmbracelet/huh"
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
	regenerateKeys := false

	_, _, err := GetKeys()
	if err == nil {
		confirmErr := huh.NewConfirm().
			Title("Secret keys already exist in the keychain, do you want to regenerate them?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&regenerateKeys).
			Run()

		if confirmErr != nil {
			return confirmErr
		}
	}

	if regenerateKeys {
		confirmAgainErr := huh.NewConfirm().
			Title("Are you sure? If you currently have secrets encrypted with the old keys you will not be able to decrypt them").
			Affirmative("Yes!").
			Negative("No.").
			Value(&regenerateKeys).
			Run()

		if confirmAgainErr != nil {
			return confirmAgainErr
		}
	}

	if !regenerateKeys {
		return nil
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
