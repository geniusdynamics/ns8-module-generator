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
	replacers := map[string]string{
		"kickstart": "nginx",
	}
	fmt.Println(replacers)
	err = parser.ReplaceInAllFiles("output", replacers)
	if err != nil {
		fmt.Printf("Error while replacing contents %v", err)
	}

}
