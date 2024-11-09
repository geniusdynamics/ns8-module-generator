package processors

import (
	"fmt"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
	"ns8-module-generator/utils"
	"strings"
)

var APP_NAME = utils.AppName

func ProcessNs8Module() {
	// Create a output Directory
	// Then do an initial commit
	err := CopyDirectory()
	if err != nil {
		fmt.Printf("error while copying directory: %v", err)
	}
	// Commit Initial Files
	err = git.GitCommitFiles("Initial commit")
	if err != nil {
		fmt.Printf("error occurred while commiting files: %s", err)
	}

	parser.DockerComposeParser(utils.DockerComposePath)
	err = ProcessBuildImage()
	if err != nil {
		fmt.Printf("error while processing build image: %v", err)
	}

	err = ReplaceAllKickstart(utils.AppName)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	GenerateMainService()
	CleanUpKickstartFiles()

	// Do git things
	if strings.ToLower(utils.AppGitInit) == "yes" {
		git.InitilaizeGitClient()

		// Push to git Online
		err = git.CreateRepository()
		if err != nil {
			fmt.Printf("An error occurred while creating repo online: %v \n", err)
			return
		}

		err = git.GitPushToRemote()
		if err != nil {
			fmt.Printf("An error occurred while pshing online: %v \n", err)
			return
		}

		fmt.Print(
			"Your app has been successfully generated. Test and see if it works as expected. Happy hacking \n",
		)
		fmt.Print("Made with ❤️  by Genius Dynamics \n")

	}
}
