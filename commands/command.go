package commands

type (
	errMsg error
)

func InputPrompts() {
	PickFile()
	InputAppName()
	InputOutputDirPath()
	InputGithubOrganizationName()
	InputGithubToken()
}
