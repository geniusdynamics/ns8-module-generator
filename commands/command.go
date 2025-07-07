package commands

import (
	"ns8-module-generator/config"
	"strings"
)

type (
	errMsg error
)

func InputPrompts(cfg *config.Config) error {
	dockerComposeSource, err := InputDockerComposeSource()
	if err != nil {
		return err
	}

	if strings.ToLower(dockerComposeSource) == "remote" {
		cfg.IsRemoteDockerCompose = true
		dockerComposeUrl, err := InputDockerComposeUrl()
		if err != nil {
			return err
		}
		cfg.DockerComposePath = dockerComposeUrl
	} else {
		cfg.IsRemoteDockerCompose = false
		file, err := PickFile()
		if err != nil {
			return err
		}
		cfg.DockerComposePath = file
	}

	appName, err := InputAppName()
	if err != nil {
		return err
	}
	cfg.AppName = strings.Join(strings.Split(appName, " "), "")

	outputDir, err := InputOutputDirPath()
	if err != nil {
		return err
	}
	cfg.OutputDir = outputDir

	gitApp, err := InputAppGitInit()
	if err != nil {
		return err
	}
	orgName, err := InputGithubOrganizationName()
	if err != nil {
		return err
	}
	cfg.GithubOrganizationName = orgName
	userName, err := InputGithubUsername()
	if err != nil {
		return err
	}
	cfg.GithubUsername = userName

	cfg.AppGitInit = strings.ToLower(gitApp) == "yes"
	if cfg.AppGitInit {
		authMethod, err := InputGitAuthMethod()
		if err != nil {
			return err
		}
		cfg.GitAuthMethod = authMethod

		if strings.ToLower(cfg.GitAuthMethod) == "token" {
			token, err := InputGithubToken()
			if err != nil {
				return err
			}
			cfg.GithubToken = token
		}
	}
	return nil
}
