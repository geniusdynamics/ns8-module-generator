package main

import (
	"ns8-module-generator/commands"
	"ns8-module-generator/processors"
	"ns8-module-generator/utils"
)

func main() {
	commands.InputPrompts()
	for utils.DockerComposePath == "" {
		commands.InputPrompts()
	}
	processors.ProcessNs8Module()
}
