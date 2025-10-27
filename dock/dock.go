package dock

import (
	"fmt"
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
	return "Add " + name + " to Dock"
}

func SetupDock(config dock.Dock, dryRun bool) {
	const delCmdPrefix = "defaults delete com.apple.dock "
	appDelErr := utils.RunCommand(delCmdPrefix+"persistent-apps", "Clear persistent apps", dryRun)
	utils.PrintError(appDelErr)

	otherDelErr := utils.RunCommand(delCmdPrefix+"persistent-others", "Clear persistent others", dryRun)
	utils.PrintError(otherDelErr)

	cmd := fmt.Sprintf("defaults write com.apple.dock show-recents -bool %t", config.ShowRecents)
	utils.RunCommand(cmd, fmt.Sprintf("Set show recents to %t", config.ShowRecents), dryRun)

	if config.Apps != nil {
		for _, path := range *config.Apps {
			var cmd string
			cmdPrefix := "defaults write com.apple.dock persistent-apps -array-add "

			if isSpacer(path) {
				cmd = fmt.Sprintf(`%s "{"tile-type"="%s-tile";}"`, cmdPrefix, path)
			} else {
				cmd = fmt.Sprintf(`%s "%s"`, cmdPrefix, appXml(path))
			}

			cmdErr := utils.RunCommand(cmd, getInfoMsg(path), dryRun)
			utils.PrintError(cmdErr)
		}

	}

	if config.Folders != nil {
		for _, path := range *config.Folders {
			cmd := fmt.Sprintf(`defaults write com.apple.dock persistent-others -array-add "%s"`, folderXml(path))
			cmdErr := utils.RunCommand(cmd, getInfoMsg(path), dryRun)
			utils.PrintError(cmdErr)
		}
	}
	utils.RunCommand("killall Dock", "Restart Dock to apply changes", dryRun)
}
