package main

import (
	"flag"
	"fmt"
	"io"
	"ns8-module-generator/commands"
	"ns8-module-generator/config"
	"ns8-module-generator/http"
	"ns8-module-generator/processors"
	"os"
)

func main() {
	cfg := config.New()
	config.Cfg = cfg

	configPath := flag.String("config", "config.yaml", "Path to the configuration file (YAML)")
	flag.Parse()

	// Check if Template Dir exists
	_, err := os.Stat(cfg.TemplateDir)
	if os.IsNotExist(err) {
		err = http.DownloadTemplate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading template: %v\n", err)
			os.Exit(1)
		}
	}

	// Try to load configuration from config.yaml
	appConfig, err := config.LoadAppConfig(*configPath)
	if err == nil {
		// If config.yaml loaded successfully, use its values
		cfg.DockerComposePath = appConfig.DockerComposePath
		cfg.AppName = appConfig.AppName
		cfg.OutputDir = appConfig.OutputDir
		cfg.AppGitInit = appConfig.AppGitInit
		cfg.GithubOrganizationName = appConfig.GithubOrganizationName
		cfg.GithubUsername = appConfig.GithubUsername
		cfg.GithubToken = appConfig.GithubToken
		cfg.GitAuthMethod = appConfig.GitAuthMethod
		cfg.IsRemoteDockerCompose = appConfig.IsRemoteDockerCompose
	} else if !os.IsNotExist(err) {
		// If config.yaml exists but there was an error loading it (e.g., malformed YAML)
		fmt.Fprintf(os.Stderr, "Error loading config.yaml: %v\n", err)
		os.Exit(1)
	} else {
		// config.yaml does not exist, fall back to interactive prompts
		err = commands.InputPrompts(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during input prompts: %v\n", err)
			os.Exit(1)
		}
	}

	var composeFileContent []byte
	if cfg.IsRemoteDockerCompose {
		// Download the file from the URL to a temporary file
		tempFile, err := os.CreateTemp("", "remote-docker-compose-*.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temporary file for remote Docker Compose: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(tempFile.Name()) // Clean up the temporary file
		defer tempFile.Close()

		err = http.DownloadFile(cfg.DockerComposePath, tempFile.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading remote Docker Compose file: %v\n", err)
			os.Exit(1)
		}
		composeFileContent, err = io.ReadAll(tempFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading downloaded Docker Compose content: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Open the local file and read the content
		composeFile, err := os.Open(cfg.DockerComposePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening local Docker Compose file: %v\n", err)
			os.Exit(1)
		}
		defer composeFile.Close()
		composeFileContent, err = io.ReadAll(composeFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading local Docker Compose content: %v\n", err)
			os.Exit(1)
		}
	}

	err = processors.ProcessNs8Module(cfg, composeFileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing NS8 module: %v\n", err)
		os.Exit(1)
	}
}
