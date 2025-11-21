package dock

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/quintisimo/macfigure/gen/dock"
	"github.com/quintisimo/macfigure/utils"
)

func appXml(path string) string {
	return fmt.Sprintf(`
		<dict>
        <key>tile-data</key>
        <dict>
            <key>file-data</key>
            <dict>
                <key>_CFURLString</key>
                <string>%s</string>
                <key>_CFURLStringType</key>
                <integer>0</integer>
            </dict>
        </dict>
    </dict>
	`, path)
}

func folderXml(path string) string {
	return fmt.Sprintf(`
		<dict>
        <key>tile-data</key>
        <dict>
            <key>arrangement</key>
            <integer>1</integer>
            <key>displayas</key>
            <integer>0</integer>
            <key>file-data</key>
            <dict>
                <key>_CFURLString</key>
                <string>file://%s</string>
                <key>_CFURLStringType</key>
                <integer>15</integer>
            </dict>
            <key>file-type</key>
            <integer>2</integer>
            <key>showas</key>
            <integer>0</integer>
        </dict>
        <key>tile-type</key>
        <string>directory-tile</string>
    </dict>
	`, path)
}

func isSpacer(path string) bool {
	return path == "spacer" || path == "small-spacer"
}

func getInfoMsg(path string) string {
	name := path
	if !isSpacer(path) {
		name = path[strings.LastIndex(path, "/") : len(path)-1]
	}
	return fmt.Sprintf("Add %s to Dock", name)
}

func updateDockItems[I any](items []I, addCmd string, rmCmd string, clrMsg string, logger *slog.Logger, dryRun bool) {
	delErr := utils.RunCommand(rmCmd, clrMsg, logger, dryRun)
	utils.PrintError(delErr, logger)

	if utils.SliceHasItems(items) {
		for _, path := range items {
			var cmd string
			cmdPrefix := fmt.Sprintf("%s -array-add", addCmd)
			path := fmt.Sprintf("%v", path)

			if isSpacer(path) {
				cmd = fmt.Sprintf(`%s "{"tile-type"="%s-tile";}"`, cmdPrefix, path)
			} else {
				xml := folderXml(path)
				if strings.HasSuffix(path, ".app") {
					xml = appXml(path)
				}

				cmd = fmt.Sprintf(`%s "%s"`, cmdPrefix, xml)
			}

			cmdErr := utils.RunCommand(cmd, getInfoMsg(path), logger, dryRun)
			utils.PrintError(cmdErr, logger)
		}

	}
}

func SetupDock(config dock.Dock, logger *slog.Logger, dryRun bool) {
	const addCmd = "defaults write com.apple.dock"
	const rmCmd = "defaults delete com.apple.dock"

	const appsCmd = "persistent-apps"
	appsAddCmd := fmt.Sprintf("%s %s", addCmd, appsCmd)
	appsRmCmd := fmt.Sprintf("%s %s", rmCmd, appsCmd)
	updateDockItems(config.Apps, appsAddCmd, appsRmCmd, "Clear persistent apps", logger, dryRun)

	const folderCmd = "persistent-others"
	foldersAddCmd := fmt.Sprintf("%s %s", addCmd, folderCmd)
	foldersRmCmd := fmt.Sprintf("%s %s", rmCmd, folderCmd)
	updateDockItems(config.Folders, foldersAddCmd, foldersRmCmd, "Clear persistent others", logger, dryRun)

	utils.WriteConfig(reflect.ValueOf(config), "com.apple.dock", addCmd, rmCmd, logger, dryRun)

	utils.RunCommand("killall Dock", "Restart Dock to apply changes", logger, dryRun)
}
