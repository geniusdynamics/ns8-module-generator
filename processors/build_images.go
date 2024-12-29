package processors

import (
	"fmt"
	"ns8-module-generator/formatters"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
	"ns8-module-generator/utils"
	"strings"
)

func ProcessBuildImage() error {
	images := formatters.GetImagesWithRepository()
	imageName := utils.AppName
	replacers := map[string]string{
		"{{ IMAGE_NAME }}":   imageName,
		"{{ GITHUB_OWNER }}": "geniusdynamics",
		"{{ IMAGES }}":       strings.Join(images, " "),
	}
	filePath := utils.OutputDir + "/build-images.sh"
	err := parser.SearchFileAndReplaceContent(filePath, replacers)
	if err != nil {
		return fmt.Errorf("error while replacing content in the file: %v", err)
	}
	// Add git file
	err = git.GitAddFile(filePath)
	if err != nil {
		return fmt.Errorf("An error occurred while adding to git: %s", err)
	}
	err = git.GitCommitFiles("feat(build-images): added images needed")
	if err != nil {
		return fmt.Errorf("An error occurred while committing: %s", err)
	}
	return nil
}
