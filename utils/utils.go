package utils

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
)

func PrintError(err error) {
	if err != nil {
		fmt.Println("[Error]", err.Error())
	}
}

func DryRunInfo(info ...string) {
	fmt.Println("[Dry Run]", info)
}

func SliceHasItems[I any](slice []I) bool {
	return len(slice) > 0
}

func GetConfigPath() string {
	homeDir := os.Getenv("XDG_CONFIG_HOME")
	if homeDir == "" {
		homeDir = fmt.Sprintf("%s/.config", os.Getenv("HOME"))
	}
	return fmt.Sprintf("%s/macfigure/config.pkl", homeDir)
}

func RunCommand(cmd string, info string, dryRun bool) error {
	if !dryRun {
		fmt.Println(info)
		command := exec.Command(cmd)
		if err := command.Run(); err != nil {
			return err
		}
	} else {
		DryRunInfo("Command:", cmd)
	}
	return nil
}

func getPropertyTypeAndValue(value reflect.Value, fieldName string) (v string, err error) {
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("-string %s", value.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("-int %d", value.Int()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("-float %f", value.Float()), nil
	case reflect.Bool:
		return fmt.Sprintf("-bool %t", value.Bool()), nil
	default:
		return "", fmt.Errorf("Unsupported type for field %s", fieldName)
	}
}

func WriteConfig(config reflect.Value, domain string, addCmd string, rmCmd string, dryRun bool) {
	for i := 0; i < config.NumField(); i++ {
		fieldName := config.Type().Field(i).Tag.Get("pkl")

		if fieldName != "apps" && fieldName != "folders" {
			field := config.Field(i)
			value := ""

			if !field.IsNil() {
				strValue, err := getPropertyTypeAndValue(field.Elem(), fieldName)
				PrintError(err)
				value = strValue
			}

			var cmd string
			var msg string
			if value != "" {
				cmd = fmt.Sprintf("%s %s %s", addCmd, fieldName, value)
				msg = "Adding"
			} else {
				cmd = fmt.Sprintf("%s %s", rmCmd, fieldName)
				msg = "Deleting"
			}

			cmdErr := RunCommand(cmd, fmt.Sprintf("%s %s %s", msg, domain, fieldName), dryRun)
			PrintError(cmdErr)
		}
	}
}
