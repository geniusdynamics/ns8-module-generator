package processors

import (
	"fmt"
	"ns8-module-generator/formatters"
	"ns8-module-generator/parser"
	"strings"
)

var OutputDirectory = "output"

func ProcessBuildImage() error {
	images := formatters.GetImagesWithRepository()
	replacers := map[string]string{
		"{{ IMAGE_NAME }}":   "nginx",
		"{{ GITHUB_OWNER }}": "geniusdynamics",
		"{{ IMAGES }}":       strings.Join(images, " "),
	}
	filePath := OutputDir + "build-images.sh"
	err := parser.SearchFileAndReplaceContent(filePath, replacers)
	if err != nil {
		return fmt.Errorf("error while replacing content in the file: %v", err)
	}
	return nil
}
