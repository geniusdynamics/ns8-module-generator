package main

import (
	"fmt"
	"ns8-module-generator/commands"
	"ns8-module-generator/config"
	"ns8-module-generator/http"
	"ns8-module-generator/processors"
	"os"
)

func main() {
	cfg := config.New()
	config.Cfg = cfg

	_, err := os.Stat(cfg.TemplateDir)
	// Check if Template Dir exists
	if os.IsNotExist(err) {
		err = http.DownloadTemplate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading template: %v\n", err)
			os.Exit(1)
		}
	}

	err = commands.InputPrompts(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during input prompts: %v\n", err)
		os.Exit(1)
	}

	err = processors.ProcessNs8Module(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing NS8 module: %v\n", err)
		os.Exit(1)
	}
}
