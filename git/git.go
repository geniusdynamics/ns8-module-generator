package git

import "github.com/google/go-github/v66/github"

func InitilaizeGitClient() {
	client := github.NewClient(nil)
}
