package config

import "github.com/go-git/go-git/v5"

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

	GitWorkTree  *git.Worktree
	GitLocalRepo *git.Repository
}

func New() *Config {
	return &Config{
		TemplateDir:    "template",
		TemplateZipURL: "https://github.com/geniusdynamics/ns8-generator-module-template/archive/refs/tags/v0.0.1.zip",
	}
}

var Cfg *Config
