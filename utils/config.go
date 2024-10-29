package utils

import "strings"

var (
	OutputDir         string
	TemplateDir       string
	DockerComposePath string
	AppName           string
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
}
