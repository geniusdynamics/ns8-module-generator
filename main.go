package main

import (
	"ns8-module-generator/commands"
	"ns8-module-generator/http"
	"ns8-module-generator/processors"
	"ns8-module-generator/utils"
	"os"
)

func main() {
	// Set utils
	utils.SetTemplateZipUrl("")

	_, err := os.Stat(utils.TemplateDir)
	// Check if Template Dir exists
	if os.IsNotExist(err) {
		http.DownloadTemplate()
	}
	commands.InputPrompts()
	for utils.DockerComposePath == "" {
		commands.InputPrompts()
	}
	processors.ProcessNs8Module()
}
