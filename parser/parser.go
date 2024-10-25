package parser

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Compose struct {
	Services map[string]Service     `yaml:"services"`
	Volumes  map[string]interface{} `yaml:"volumes"`
	Networks map[string]interface{} `yaml:"networks"`
}
type Images struct {
	// An array of images
	Images []string
}

type Service struct {
	Name        string
	Image       string      `yaml:"image"`
	Environment []string    `yaml:"environment"`
	DependsOn   interface{} `yaml:"depends_on,omitempty"`
	Volumes     []string    `yaml:"volumes,omitempty"`
}

// Global variable to hold the images
var (
	images   Images
	services []Service
)

func appendImage(image string) {
	images.Images = append(images.Images, image)
}

func appendServices(service Service) {
	services = append(services, service)
}

func GetServices() *[]Service {
	return &services
}

func ParseDockerCompose(
	filePath string,
) (map[string]Service, map[string]interface{}, map[string]interface{},
	error,
) {
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

	// Extract and store images globally
	for _, service := range compose.Services {
		images.Images = append(images.Images, service.Image)
	}
	/*
		Return the services, volumes and networks
	*/
	return compose.Services, compose.Volumes, compose.Networks, nil
}

// ParseServiceContents /*
func ParseServiceContents(services map[string]Service) {
	/*
	* Loop through the services
	 */
	for name, service := range services {
		fmt.Printf("Parsing through service: %v \n", name)
		service.Name = name
		appendServices(service)
	}
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

// GetImages Get all images and return them
func GetImages() []string {
	return images.Images
}

// Deal with Docker Compose file

func DockerComposeParser(filename string) {
	services, volumes, networks, e := ParseDockerCompose(filename)
	if e != nil {
		fmt.Printf("Error while parsing docker-compose file: %v", e)
	}
	ParseServiceContents(services)
	ParseVolumeContents(volumes)
	ParseNetworkContents(networks)
}
