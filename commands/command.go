package commands

import (
	"ns8-module-generator/utils"
	"strings"
)

type (
	errMsg error
)

func InputPrompts() {
	PickFile()
	InputAppName()
	InputOutputDirPath()
	InputAppGitInit()
	if strings.ToLower(utils.AppGitInit) == "yes" {
		InputGithubOrganizationName()
		InputGithubUsername()
		InputGithubToken()
	}
}
