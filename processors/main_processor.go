package processors

import (
	"fmt"
	"ns8-module-generator/parser"
	"ns8-module-generator/utils"
)

var APP_NAME = "nginx"

func ProcessNs8Module() {
	utils.SetOutputDir("output")
	utils.SetTemplateDir("template")
	// Create a temp Directory
	err := CopyDirectory()
	if err != nil {
		fmt.Printf("error while copying directory: %v", err)
	}
	parser.DockerComposeParser("./tests/docker-compose.yaml")
	err = ProcessBuildImage()
	if err != nil {
		fmt.Printf("error while processing build image: %v", err)
	}

	err = ReplaceAllKickstart("nginx1")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	GenerateMainService()
}
