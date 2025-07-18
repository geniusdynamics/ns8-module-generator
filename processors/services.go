package processors

import (
	"bufio"
	"fmt"
	"ns8-module-generator/config"
	"ns8-module-generator/formatters"
	"ns8-module-generator/generators"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
	"os"
	"strings"
)

type ServiceFile struct {
	// Unit
	Description string `json:"description,omitempty"`
	Requires    string `json:"requires,omitempty"`
	Before      string `json:"before,omitempty"`
	// Service
	ExecStartPre    string `json:"exec_start_pre,omitempty"`
	ExecStart       string `json:"exec_start,omitempty"`
	ExecStop        string `json:"exec_stop,omitempty"`
	ExecStopPost    string `json:"exec_stop_post,omitempty"`
	PIDFile         string `json:"pid_file,omitempty"`
	Restart         string `json:"restart,omitempty"`
	TimeoutStopSec  string `json:"timeout_stop_sec,omitempty"`
	EnvironmentFile string `json:"environment_file,omitempty"`
	Type            string `json:"type,omitempty"`
	PublishPort     string `json:"publish_port,omitempty"`
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

// readServiceFile function  
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

// GenerateMainService function  
func GenerateMainService() {
	// Read the main Service file

	// Read The Service file
	service, e := readServiceFileContents(
		config.Cfg.OutputDir + "/imageroot/systemd/user/kickstart.service",
	)
	if e != nil {
		fmt.Printf("An error occurred: %v", e)
	}
	images := formatters.GetImagesCompatibleServiceNames()
	fmt.Printf("All Images: %v", images)

	// Generete required services
	requiredServices := generateRequiredServices(images)
	fmt.Printf("All required Services: %v", requiredServices)

	allServices := generateRequiredServices(images)
	// Replacers
	replacers := map[string]string{
		"{{ SERVICE_NAME }}":      config.Cfg.AppName,
		"{{ REQUIRED_SERVICES }}": allServices,
		"{{ BEFORE_SERVICES }}":   allServices,
	}
	formattedContent := formatters.ReplacePlaceHolders(service, replacers)
	print(formattedContent)

	err := writeServiceFile(formattedContent, config.Cfg.AppName+".service")
	if err != nil {
		fmt.Printf("An error occurred while writing to service: %v", err)
		return
	}
	GenerateServicesFiles(allServices)
}

func GenerateServicesFiles(allServices string) {
	serviceContent, e := readServiceFileContents(
		config.Cfg.OutputDir + "/imageroot/systemd/user/kickstart-app.service",
	)
	if e != nil {
		fmt.Printf("An error occurred reading kickstart-app.service: %v", e)
		return
	}
	for _, service := range *parser.GetServices() {
		var envPath string
		if generators.IsCommonDatabaseImage(service.Image) {
			envPath = config.Cfg.OutputDir + "/imageroot/actions/create-module/10configure_environment_vars"
		} else {
			envPath = config.Cfg.OutputDir + "/imageroot/actions/configure-module/10configure_environment_vars"
		}

		// Generate environment file contents
		env, err := generators.GenerateEnvFileContents(service.Name, service.ParsedEnvironment, envPath)
		cleanEnv := strings.TrimSpace(strings.TrimPrefix(env, "--env-file"))

		if generators.IsCommonDatabaseImage(service.Image) {
			restore, backup, clean := generators.GenerateBackupRestore(
				formatters.ImageNameWithSuffix(service.Image),
				service.Name,
				cleanEnv,
				service.Name,
			)
			fmt.Print("Restore: ", restore)
			fmt.Print("Backup:", backup)
			fmt.Print("Clean: ", clean)
			generators.WriteToFile(config.Cfg.OutputDir+"/imageroot/actions/restore-module/40restore_database", restore, "feat: added Restore")
			generators.WriteToFile(config.Cfg.OutputDir+"/imageroot/bin/module-dump-state", backup, "feat: added backup")
			generators.WriteToFile(config.Cfg.OutputDir+"/imageroot/bin/module-cleanup-state", clean, "feat: added clean up")
		}
		if err != nil {
			fmt.Printf("Failed to generate env file for service %s: %v\n", service.Name, err)
			return
		}
		replacers := map[string]string{
			"{{ SERVICE_NAME }}":      service.Name + "-app",
			"{{ MAIN_SERVICE_NAME }}": config.Cfg.AppName,
			"{{ IMAGE_NAME }}":        formatters.ImageNameWithSuffix(service.Image),
			"{{ OTHER_COMMANDS }}":    "",
			"{{ VOLUMES }}":           generators.GenerateNS8VolumeFlags(service.ParsedVolumes),
			"{{ AFTER_SERVICES }}": generators.GenerateNS8AfterServices(
				service.DependsOn,
				allServices,
				config.Cfg.AppName+".service",
			),
			"{{ BINDS_TO_SERVICES }}": config.Cfg.AppName + ".service",
			"{{ ENV_FILES }}":         env,
		}
		formattedServiceContent := formatters.ReplacePlaceHolders(serviceContent, replacers)
		// print(formattedServiceContent)

		// Generate Get Configuration content
		err = generators.GenerateGetConfigurationContent(
			service.Name,
			service.ParsedEnvironment,
			config.Cfg.OutputDir+"/imageroot/actions/get-configuration/20read",
		)
		if err != nil {
			fmt.Printf("An error occured while writing get configuration content: %v\n", err)
			return
		}
		// Save the service file
		e = writeServiceFile(formattedServiceContent, service.Name+"-app.service")
		if e != nil {
			fmt.Errorf(
				"An error occurred while saving service file: %v \n",
				service.Name+"-app.service",
			)
		}

	}
	// Add Json dump at the end
	err := generators.AddJsonDump(config.Cfg.OutputDir + "/imageroot/actions/get-configuration/20read")
	if err != nil {
		fmt.Printf("An error occurred adding json dump: %v\n", err)
		return
	}
}

func writeServiceFile(content string, fileName string) error {
	filePath := config.Cfg.OutputDir + "/imageroot/systemd/user/" + fileName
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	err = git.GitAddFile(filePath)
	if err != nil {
		return err
	}
	err = git.GitCommitFiles("feat(service): added " + fileName)
	if err != nil {
		return err
	}
	return nil
}

func generateRequiredServices(images []string) string {
	services := ""
	for index := range images {
		services += images[index] + "-app.service "
	}
	return services
}
