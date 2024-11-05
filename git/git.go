package git

import (
	"context"
	"fmt"
	"ns8-module-generator/utils"

	"github.com/google/go-github/v66/github"
)

var client *github.Client

func InitilaizeGitClient() {
	// Get Github token
	token := utils.GithubToken
	// Create new git client
	client = github.NewClient(nil).WithAuthToken(token)
}

func CreateRepository() error {
	repo := &github.Repository{
		Name:    github.String(""),
		Private: github.Bool(false),
	}
	ctx := context.Background()
	// Return response err and repo details
	repo, _, err := client.Repositories.Create(ctx, utils.GithubOrganizationName, repo)
	if err != nil {
		return fmt.Errorf(
			"An error occurred while occurred while creating the Repository: %v \n",
			err,
		)
	}
	fmt.Printf("The Git URL: %s", repo.GetGitURL())
	utils.SetGitRemoteUrl(repo.GetGitURL())
	return nil
}
