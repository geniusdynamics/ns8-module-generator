package parser

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Compose struct {
	Services map[string]interface{} `yaml:"services"`
	Volumes  map[string]interface{} `yaml:"volumes"`
	Networks map[string]interface{} `yaml:"networks"`
}

func ParseDockerCompose(filePath string) (map[string]interface{}, map[string]interface{}, map[string]interface{},
	error) {
	/*
		Open the file and read the content
	*/
	composeFile, e := os.Open(filePath)
	/*
		Check if there is an error
	*/
	if e != nil {
		return nil, nil, nil, e
	}
	/*
		Defer the closing of the file
	*/
	defer func(composeFile *os.File) {
		err := composeFile.Close()
		if err != nil {
			panic(err)
		}
	}(composeFile)
	byteValue, _ := io.ReadAll(composeFile)
	var compose Compose
	e = yaml.Unmarshal(byteValue, &compose)
	// Check if there is an error
	if e != nil {
		return nil, nil, nil, e
	}
	/*
		Return the services, volumes and networks
	*/
	return compose.Services, compose.Volumes, compose.Networks, nil
}

// ParseServiceContents /*
func ParseServiceContents(services map[string]interface{}) {
	/*
	* Loop through the services
	 */
	for name, value := range services {
		println("Service: ", name)
		for key1, value1 := range value.(map[string]interface{}) {
			println("Key: ", key1)
			fmt.Printf("Value: %v \n", value1)
		}
	}
}

func getImageName() {

}

func ParseVolumeContents(volume map[string]interface{}) {
	for key, value := range volume {
		println("Volume: ", key)
		fmt.Printf("Value: %v \n", value)
	}
}

func serviceImages() {

}
func ParseNetworkContents(network map[string]interface{}) {
	for key, value := range network {
		println("Network: ", key)
		fmt.Printf("Value: %v \n", value)
	}
}
