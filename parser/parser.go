package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

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
	Name            string
	Image           string      `yaml:"image"`
	Environment     yaml.Node   `yaml:"environment"`
	ParsedEnvironment map[string]string
	DependsOn       interface{} `yaml:"depends_on,omitempty"`
	Volumes         yaml.Node   `yaml:"volumes,omitempty"` // Change to yaml.Node
	ParsedVolumes   []map[string]string // New field for parsed volumes
}

func (s *Service) UnmarshalYAML(node *yaml.Node) error {
	type rawService Service
	if err := node.Decode((*rawService)(s)); err != nil {
		return err
	}

	s.ParsedEnvironment = make(map[string]string)
	if s.Environment.Kind == yaml.SequenceNode {
		for _, envNode := range s.Environment.Content {
			if envNode.Kind == yaml.ScalarNode {
				parts := strings.SplitN(envNode.Value, "=", 2)
				if len(parts) == 2 {
					s.ParsedEnvironment[parts[0]] = parts[1]
				} else {
					s.ParsedEnvironment[parts[0]] = "" // Handle cases like - VAR_NAME
				}
			}
		}
	} else if s.Environment.Kind == yaml.MappingNode {
		for i := 0; i < len(s.Environment.Content); i += 2 {
			keyNode := s.Environment.Content[i]
			valueNode := s.Environment.Content[i+1]
			if keyNode.Kind == yaml.ScalarNode && valueNode.Kind == yaml.ScalarNode {
				s.ParsedEnvironment[keyNode.Value] = valueNode.Value
			}
		}
	}

	s.ParsedVolumes = []map[string]string{}
	if s.Volumes.Kind == yaml.SequenceNode {
		for _, volNode := range s.Volumes.Content {
			if volNode.Kind == yaml.ScalarNode {
				// Short syntax: /host/path:/container/path
				parts := strings.SplitN(volNode.Value, ":", 2)
				if len(parts) == 2 {
					s.ParsedVolumes = append(s.ParsedVolumes, map[string]string{"source": parts[0], "target": parts[1]})
				} else {
					s.ParsedVolumes = append(s.ParsedVolumes, map[string]string{"source": parts[0], "target": parts[0]})
				}
			} else if volNode.Kind == yaml.MappingNode {
				// Long syntax: type: bind, source: /host, target: /container
				volumeMap := make(map[string]string)
				for i := 0; i < len(volNode.Content); i += 2 {
					key := volNode.Content[i].Value
					value := volNode.Content[i+1].Value
					volumeMap[key] = value
				}
				s.ParsedVolumes = append(s.ParsedVolumes, volumeMap)
			}
		}
	}

	return nil
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
		fmt.Printf("Parsed Environment for %s: %+v\n", name, service.ParsedEnvironment)
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
