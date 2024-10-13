package processors

import (
	"bufio"
	"fmt"
	"ns8-module-generator/formatters"
	"os"
	"strings"
)

type ServiceFile struct {
	// Unit
	Description string
	Requires    string
	Before      string
	// Service
	ExecStartPre    string
	ExecStart       string
	ExecStop        string
	ExecStopPost    string
	PIDFile         string
	Restart         string
	TimeoutStopSec  string
	EnvironmentFile string
	Type            string
	PublishPort     string
}

//	func assignField(service *ServiceFile, line string) {
//		fields := []string{
//			"Description", "Requires", "Before", "ExecStartPre", "ExecStart", "ExecStop", "ExecStopPost",
//			"PIDFile", "Restart", "TimeoutStopSec", "EnvironmentFile", "Type",
//		}
//		for _, field := range fields {
//			if strings.HasPrefix(line, field+"=") {
//				service.Fields[field] = line
//				return
//			}
//		}
//	}

func readServiceFileContents(filePath string) (string, error) {
	serviceFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer serviceFile.Close()

	var contents strings.Builder
	scanner := bufio.NewScanner(serviceFile)

	var currentCommand string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Handle multi-line commands
		if strings.HasSuffix(line, "\\") {
			currentCommand += strings.TrimSuffix(line, "\\") + " "
			continue
		}

		// Append the last part of the multi-line command
		if currentCommand != "" {
			currentCommand += line
			line = currentCommand
			currentCommand = ""
		}

		// Append the line to the contents builder
		contents.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return contents.String(), nil
}

func readServiceFile(filePath string) (*ServiceFile, error) {
	serviceFile, e := os.Open(filePath)
	if e != nil {
		return nil, e
	}
	// Close The file
	defer serviceFile.Close()
	// Read All the contents in the service file
	// Service
	service := &ServiceFile{}

	scanner := bufio.NewScanner(serviceFile)
	var currentCommand string
	for scanner.Scan() {
		// Read The Lines
		line := strings.TrimSpace(scanner.Text())
		// Check if The line Contains comments
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Handle Multi Lines
		if strings.HasSuffix(line, "\\") {
			currentCommand += strings.TrimSuffix(line, "\\") + ""
			// Skip if has backlashes at end
			continue
		}
		// Apend Last Part
		if currentCommand != "" {
			currentCommand += line
			line = currentCommand
			currentCommand = ""
		}

		// Check if Has Other Components
		// Now we parse individual service file fields
		if strings.HasPrefix(line, "Description=") {
			service.Description = line
		} else if strings.HasPrefix(line, "Requires=") {
			service.Requires = line
		} else if strings.HasPrefix(line, "Before=") {
			service.Before = line
		} else if strings.HasPrefix(line, "ExecStartPre=") {
			service.ExecStartPre = line
		} else if strings.HasPrefix(line, "ExecStart=") {
			service.ExecStart = line
		} else if strings.HasPrefix(line, "ExecStop=") {
			service.ExecStop = line
		} else if strings.HasPrefix(line, "ExecStopPost=") {
			service.ExecStopPost = line
		} else if strings.HasPrefix(line, "PIDFile=") {
			service.PIDFile = line
		} else if strings.HasPrefix(line, "Restart=") {
			service.Restart = line
		} else if strings.HasPrefix(line, "TimeoutStopSec=") {
			service.TimeoutStopSec = line
		} else if strings.HasPrefix(line, "EnvironmentFile=") {
			service.EnvironmentFile = line
		} else if strings.HasPrefix(line, "Type=") {
			service.Type = line
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// Return The service
	return service, nil
}

func GenerateMainService() {
	// appName := APP_NAME
	// Read the main Service file

	// Read The Service file
	// service, e := readServiceFileContents("./template/imageroot/systemd/user/kickstart.service")
	// if e != nil {
	// 	fmt.Printf("An error occurred: %v", e)
	// }
	images := formatters.GetImagesCompatibleServiceNames()
	fmt.Printf("All Images: %v", images)
	// Replacers
	// replacers := map[string]string{}
}
