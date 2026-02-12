package utils

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/charmbracelet/log"
)

func DryRunInfo(info string, logger *log.Logger) {
	logger.With("type", "dry-run").Info(info)
}

func SliceHasItems[I any](slice []I) bool {
	return len(slice) > 0
}

func RunCommand(cmd string, info string, logger *log.Logger, dryRun bool) error {
	if !dryRun {
		log.Info(info)
		command := exec.Command(cmd)
		if err := command.Run(); err != nil {
			return err
		}
	} else {
		DryRunInfo(fmt.Sprintf("Running %s", cmd), logger)
	}
	return nil
}

func CopyFile(source string, target string, logger *log.Logger, dryRun bool) error {
	if !dryRun {
		contents, readErr := os.ReadFile(source)
		if readErr != nil {
			return readErr
		}

		file, createErr := os.Create(target)
		if createErr != nil {
			return createErr
		}
		defer file.Close()

		_, writeErr := file.Write(contents)
		if writeErr != nil {
			return writeErr
		}
	} else {
		DryRunInfo(fmt.Sprintf("Creating %s", target), logger)
	}
	return nil
}

func getPropertyTypeAndValue(value reflect.Value, fieldName string) (string, error) {
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

func WriteConfig(config reflect.Value, domain string, addCmd string, rmCmd string, logger *log.Logger, dryRun bool) error {
	for i := 0; i < config.NumField(); i++ {
		fieldName := config.Type().Field(i).Tag.Get("pkl")

		if fieldName != "apps" && fieldName != "folders" {
			field := config.Field(i)
			value := ""

			if !field.IsNil() {
				strValue, err := getPropertyTypeAndValue(field.Elem(), fieldName)
				if err != nil {
					return err
				}
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

			cmdErr := RunCommand(cmd, fmt.Sprintf("%s %s %s", msg, domain, fieldName), logger, dryRun)
			if cmdErr != nil {
				return cmdErr
			}
		}
	}
	return nil
}
