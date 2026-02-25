package utils

import (
	"fmt"
	"io"
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

		if commandErr := exec.Command(cmd).Run(); commandErr != nil {
			return commandErr
		}
	} else {
		DryRunInfo(fmt.Sprintf("Running %s", cmd), logger)
	}
	return nil
}

func ReadFile(source string, logger *log.Logger, dryRun bool) (io.Reader, error) {
	if !dryRun {
		file, openErr := os.Open(source)
		if openErr != nil {
			return nil, openErr
		}
		return file, nil
	} else {
		DryRunInfo(fmt.Sprintf("Reading %s", source), logger)
	}
	return nil, nil
}

func WriteFile(reader io.Reader, target string, logger *log.Logger, dryRun bool) error {
	if !dryRun {
		file, createErr := os.Create(target)
		if createErr != nil {
			return createErr
		}
		defer file.Close()

		if _, writeErr := io.Copy(file, reader); writeErr != nil {
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
				strValue, typeAndValueErr := getPropertyTypeAndValue(field.Elem(), fieldName)
				if typeAndValueErr != nil {
					return typeAndValueErr
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

			if cmdErr := RunCommand(cmd, fmt.Sprintf("%s %s %s", msg, domain, fieldName), logger, dryRun); cmdErr != nil {
				return cmdErr
			}
		}
	}
	return nil
}
