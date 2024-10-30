package utils

import "strings"

var (
	OutputDir         string
	TemplateDir       string
	DockerComposePath string
	AppName           string
	TemplateZipURL    string
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
