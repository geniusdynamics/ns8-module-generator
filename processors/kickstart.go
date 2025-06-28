package processors

import (
	"fmt"
	"ns8-module-generator/config"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
	"os"
)

func ReplaceAllKickstart(appName string) error {
	replacers := map[string]string{
		"kickstart": appName,
	}

	err := parser.ReplaceInAllFiles(config.Cfg.OutputDir, replacers)
	if err != nil {
		return fmt.Errorf("An error occurred: %v", err)
	}
	err = git.GitCommitFiles("refactor: replaced all kickstart names with " + config.Cfg.AppName)
	if err != nil {
		return err
	}
	return nil
}

func CleanUpKickstartFiles() {
	filePaths := []string{
		config.Cfg.OutputDir + "/imageroot/systemd/user/kickstart.service",
		config.Cfg.OutputDir + "/imageroot/systemd/user/kickstart-app.service",
	}
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("An error occurred while deleting %s : %v \n", filePath, err)
		}
		// Git add the removed file
		err = git.GitAddFile(filePath)
		if err != nil {
			fmt.Printf("An error occurred: %v", err)
		}
	}
	err := git.GitCommitFiles("Removed kickstart files")
	if err != nil {
		fmt.Printf("An eeror occurred while committing", err)
	}
}
