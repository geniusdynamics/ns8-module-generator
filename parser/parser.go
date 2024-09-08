package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Compose struct {
	Services map[string]interface{} `yaml:"services"`
	Volumes  map[string]interface{} `yaml:"volumes"`
	Networks map[string]interface{} `yaml:"networks"`
}

func parseDockerCompose(filePath string) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	composeFile, e := os.Open(filePath)
	if e != nil {
		return nil, nil, nil, e
	}
	defer func(composeFile *os.File) {
		err := composeFile.Close()
		if err != nil {
			panic(err)
		}
	}(composeFile)
	byteValue, _ := io.ReadAll(composeFile)
	var compose Compose
	e = yaml.Unmarshal(byteValue, &compose)
	if e != nil {
		return nil, nil, nil, e
	}
	return compose.Services, compose.Volumes, compose.Networks, nil
}
