package git

import (
	"context"
	"fmt"
	"strings"

	"ns8-module-generator/config"

	"github.com/google/go-github/v66/github"
)

var client *github.Client

func InitilaizeGitClient(cfg *config.Config) {
	// Get Github token
	// Create new git client
	client = github.NewClient(nil).WithAuthToken(cfg.GithubToken)
}

func CreateRepository(cfg *config.Config) error {
	name := cfg.AppName
	repo := &github.Repository{
		Name:          github.String("ns8-" + name),
		Private:       github.Bool(false),
		DefaultBranch: github.String("main"),
	}
	ctx := context.Background()
	// Return response err and repo details
	repo, _, err := client.Repositories.Create(ctx, cfg.GithubOrganizationName, repo)
	if err != nil {
		return fmt.Errorf(
			"An error occurred while occurred while creating the Repository: %v \n",
			err,
		)
	}
	// Print the repository URL
	fmt.Printf("The Git URL: %s \n", repo.GetHTMLURL())
	if strings.ToLower(config.Cfg.GitAuthMethod) == "ssh" {
		cfg.GitRemoteUrl = repo.GetSSHURL()
	} else {
		cfg.GitRemoteUrl = repo.GetHTMLURL()
	}

	return nil
}
