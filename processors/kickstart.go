package processors

import (
	"fmt"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
	"ns8-module-generator/utils"
	"os"
)

func ReplaceAllKickstart(appName string) error {
	replacers := map[string]string{
		"kickstart": appName,
	}

	err := parser.ReplaceInAllFiles(utils.OutputDir, replacers)
	if err != nil {
		return fmt.Errorf("An error occurred: %v", err)
	}
	err = git.GitCommitFiles("refactor: replaced all kickstart names with " + utils.AppName)
	if err != nil {
		return err
	}
	return nil
}

func CleanUpKickstartFiles() {
	filePaths := []string{
		utils.OutputDir + "/imageroot/systemd/user/kickstart.service",
		utils.OutputDir + "/imageroot/systemd/user/kickstart-app.service",
	}
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("An error occurred while deleting %s : %v \n", filePath, err)
		}
	}
}
