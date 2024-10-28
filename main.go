package main

import (
	"ns8-module-generator/commands"
	"ns8-module-generator/processors"
)

func main() {
	commands.PropmtInputs()
	processors.ProcessNs8Module()
}
