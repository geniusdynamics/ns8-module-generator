package commands

import (
	"log"
	"ns8-module-generator/config"
	"strings"
)

type (
	errMsg error
)

func InputPrompts(cfg *config.Config) {
	file, err := PickFile()
	if err != nil {
		log.Fatal(err)
	}
	cfg.DockerComposePath = file
	appName, err := InputAppName()
	if err != nil {
		log.Fatal(err)
	}
	cfg.AppName = strings.Join(strings.Split(appName, " "), "")

	outputDir, err := InputOutputDirPath()
	if err != nil {
		log.Fatal(err)
	}
	cfg.OutputDir = outputDir

	gitApp, err := InputAppGitInit()
	if err != nil {
		log.Fatal(err)
	}
	cfg.AppGitInit = strings.ToLower(gitApp) == "yes"
	if cfg.AppGitInit {
		orgName, err := InputGithubOrganizationName()
		if err != nil {
			log.Fatal(err)
		}
		cfg.GithubOrganizationName = orgName
		userName, err := InputGithubUsername()
		if err != nil {
			log.Fatal(err)
		}
		cfg.GithubUsername = userName
		token, err := InputGithubToken()
		if err != nil {
			log.Fatal(err)
		}
		cfg.GithubToken = token
	}
}
