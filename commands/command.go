package commands

import (
	"fmt"
	"ns8-module-generator/utils"
	"os"

	"github.com/manifoldco/promptui"
)

func PropmtInputs() {
	propmt := promptui.Prompt{
		Label: "Path to you docker compose",
		Validate: func(input string) error {
			info, err := os.Stat(input)
			if err != nil {
				return fmt.Errorf("The file does not exist")
			}
			if info.IsDir() {
				return fmt.Errorf("This is not a file. Its a directory")
			}
			return nil
		},
	}
	result, err := propmt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v \n", err)
	}
	fmt.Printf("Docker compose path: %q \n", result)
	utils.SetDockerComposePath(result)
}
