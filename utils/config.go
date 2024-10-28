package utils

var (
	OutputDir         string
	TemplateDir       string
	DockerComposePath string
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
