package config

import (
	"os"

	"github.com/go-git/go-git/v5"
	gopkg_yaml_v3 "gopkg.in/yaml.v3"
)

type Config struct {
	DockerComposePath      string
	AppName                string
	OutputDir              string
	AppGitInit             bool
	GithubOrganizationName string
	GithubUsername         string
	GithubToken            string
	TemplateDir            string
	TemplateZipURL         string
	GitRemoteUrl           string
	GitAuthMethod          string
	IsRemoteDockerCompose  bool

	GitWorkTree  *git.Worktree
	GitLocalRepo *git.Repository
}

func New() *Config {
	return &Config{
		TemplateDir:    "template",
		TemplateZipURL: "https://github.com/geniusdynamics/ns8-generator-module-template/archive/refs/tags/v0.0.2.zip",
	}
}

var Cfg *Config

func LoadAppConfig(filePath string) (*AppConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig
	if err := gopkg_yaml_v3.Unmarshal(data, &appConfig); err != nil {
		return nil, err
	}

	return &appConfig, nil
}
