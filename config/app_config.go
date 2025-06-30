package config

type AppConfig struct {
	DockerComposePath      string `yaml:"dockerComposePath"`
	IsRemoteDockerCompose  bool   `yaml:"isRemoteDockerCompose"`
	AppName                string `yaml:"appName"`
	OutputDir              string `yaml:"outputDir"`
	AppGitInit             bool   `yaml:"appGitInit"`
	GithubOrganizationName string `yaml:"githubOrganizationName"`
	GithubUsername         string `yaml:"githubUsername"`
	GithubToken            string `yaml:"githubToken"`
	GitAuthMethod          string `yaml:"gitAuthMethod"`
}
