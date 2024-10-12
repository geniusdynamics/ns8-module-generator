package processors

import (
	"fmt"
	"ns8-module-generator/parser"
)

func ProcessNs8Module() {
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
}
