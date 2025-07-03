package utils

import (
	"strings"

	"github.com/go-git/go-git/v5"
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

	GitWorkTree  *git.Worktree
	GitLocalRepo *git.Repository
}

func New() *Config {
	return &Config{
		TemplateDir:    "template",
		TemplateZipURL: "https://github.com/geniusdynamics/ns8-generator-module-template/archive/refs/tags/v0.0.1.zip",
	}
}

var (
	OutputDir              string
	TemplateDir            string
	DockerComposePath      string
	AppName                string
	TemplateZipURL         string
	GithubToken            string
	GithubUsername         string
	GithubOrganizationName string
	AppGitInit             string
	GitRemoteUrl           string
	GitWorkTree            *git.Worktree
	GitLocalRepo           *git.Repository
)

func SetOutputDir(dir string) {
	OutputDir = dir
}

func SetTemplateDir(dir string) {
	TemplateDir = dir
}

func SetDockerComposePath(path string) {
	DockerComposePath = path
}

func SetAppName(appName string) {
	app_name := strings.Split(appName, " ")

	AppName = strings.Join(app_name, "")
	// AppName = appName
}

func SetTemplateZipUrl(url string) {
	TemplateZipURL = url
}

func SetGithubToken(token string) {
	GithubToken = token
}

func SetGithubUsername(name string) {
	GithubUsername = name
}

func SetGithubOrganizationName(name string) {
	GithubOrganizationName = name
}

func SetAppGitInit(val string) {
	AppGitInit = val
}

func SetGitRemoteUrl(url string) {
	GitRemoteUrl = url
}

func SetGitWorkingTree(worktree *git.Worktree) {
	GitWorkTree = worktree
}

func SetGitLocalRepo(repo *git.Repository) {
	GitLocalRepo = repo
}
