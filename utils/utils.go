package utils

import (
	"fmt"
	"os/exec"
	"reflect"
)

func SliceDifference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func DryRunInfo(info string) {
	fmt.Println("[Dry Run] " + info)
}

func RunCommand(cmd string, info string, dryRun bool) {
	if !dryRun {
		fmt.Println(info)
		command := exec.Command(cmd)
		if err := command.Run(); err != nil {
			panic(err)
		}
	} else {
		DryRunInfo("[Dry Run] Command: " + cmd)
	}
}

func CovertToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", value.Bool())
	default:
		return ""
	}
}
