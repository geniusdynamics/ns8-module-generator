package parser

import (
	"fmt"
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

// ParseComposeContent parses the byte content of a Docker Compose file.
// It populates the global images and services variables.
func ParseComposeContent(content []byte) (map[string]Service, map[string]interface{}, map[string]interface{}, error) {
	var compose Compose
	err := yaml.Unmarshal(content, &compose)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while parsing docker-compose content: %w", err)
	}

	// Clear previous images and services before populating
	images.Images = nil
	services = nil

	// Extract and store images globally
	for _, service := range compose.Services {
		images.Images = append(images.Images, service.Image)
	}

	ParseServiceContents(compose.Services)
	ParseVolumeContents(compose.Volumes)
	ParseNetworkContents(compose.Networks)

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

