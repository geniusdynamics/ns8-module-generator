package processors

import (
	"fmt"
	"ns8-module-generator/parser"
	"ns8-module-generator/utils"
)

func ReplaceAllKickstart(appName string) error {
	replacers := map[string]string{
		"kickstart": appName,
	}

	err := parser.ReplaceInAllFiles(utils.OutputDir, replacers)
	if err != nil {
		return fmt.Errorf("An error occurred: %v", err)
	}
	return nil
}
