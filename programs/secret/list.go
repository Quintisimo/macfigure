package secret

import (
	"github.com/charmbracelet/huh"
	"github.com/quintisimo/macfigure/gen/secret"
)

func List(secret []secret.Secret) (string, error) {
	const defaultWidth = 20

	secretPath := ""
	secretLen := len(secret)

	if secretLen == 0 {
		return "", nil
	}

	if secretLen == 1 {
		return secret[0].Source, nil
	}

	items := make([]string, secretLen)
	for i, secretItem := range secret {
		items[i] = secretItem.Source
	}

	if runErr := huh.NewSelect[string]().
		Title("Select secret to edit").
		Options(huh.NewOptions(items...)...).
		Value(&secretPath).
		Run(); runErr != nil {
		return "", runErr
	}

	return secretPath, nil
}
