package processors

import (
	"fmt"
	"ns8-module-generator/config"
	"ns8-module-generator/git"
	"ns8-module-generator/parser"
)

func ProcessNs8Module(cfg *config.Config, composeFileContent []byte) error {
	// Create a output Directory
	// Then do an initial commit
	err := CopyDirectory()
	if err != nil {
		return fmt.Errorf("error while copying directory: %v", err)
	}
	// Commit Initial Files
	err = git.GitCommitFiles("Initial commit")
	if err != nil {
		return fmt.Errorf("error occurred while commiting files: %s", err)
	}

	_, _, _, err = parser.ParseComposeContent(composeFileContent)
	if err != nil {
		return fmt.Errorf("Error parsing Docker Compose content: %v", err)
	}

	err = ProcessBuildImage()
	if err != nil {
		return fmt.Errorf("error while processing build image: %v", err)
	}

	err = ReplaceAllKickstart(cfg.AppName)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	GenerateMainService()
	CleanUpKickstartFiles()

	// Do git things
	if cfg.AppGitInit {
		git.InitilaizeGitClient(cfg)

		// Push to git Online
        err = git.CreateRepository(cfg)
        if err != nil {
            return fmt.Errorf("An error occurred while creating repo online: %v \n", err)
        }

        err = git.GitPushToRemote()
        if err != nil {
            return fmt.Errorf("An error occurred while pshing online: %v \n", err)
        }

        fmt.Print(
            "Your app has been successfully generated. Test and see if it works as expected. Happy hacking \n",
        )
        fmt.Print("Made with ❤️  by Genius Dynamics \n")

    }
    return nil
}
