package main

import (
	"ns8-module-generator/commands"
	"ns8-module-generator/config"
	"ns8-module-generator/http"
	"ns8-module-generator/processors"
	"os"
)

func main() {
	cfg := config.New()
	config.Cfg = cfg
	// config.Cfg.SetOutputDir("output")

	_, err := os.Stat(cfg.TemplateDir)
	// Check if Template Dir exists
	if os.IsNotExist(err) {
		http.DownloadTemplate()
	}
	commands.InputPrompts(cfg)
	// for config.Cfg.DockerComposePath == "" {
	// 	commands.InputPrompts()
	// }
	processors.ProcessNs8Module(cfg)
}
