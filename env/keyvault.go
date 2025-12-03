package env

import (
	"errors"

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
