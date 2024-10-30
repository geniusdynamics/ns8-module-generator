package main

import (
	"ns8-module-generator/commands"
	"ns8-module-generator/http"
	"ns8-module-generator/utils"
	"os"
)

func main() {
	utils.SetOutputDir("output")
	utils.SetTemplateDir("template")
	// Set utils
	utils.SetTemplateZipUrl(
		"https://github.com/geniusdynamics/ns8-generator-module-template/archive/refs/tags/v0.0.1.zip",
	)

	_, err := os.Stat(utils.TemplateDir)
	// Check if Template Dir exists
	if os.IsNotExist(err) {
		http.DownloadTemplate()
	}
	commands.InputPrompts()
	// for utils.DockerComposePath == "" {
	// 	commands.InputPrompts()
	// }
	// processors.ProcessNs8Module()
}
