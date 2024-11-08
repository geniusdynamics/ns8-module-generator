package git

import (
	"fmt"
	"ns8-module-generator/utils"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func InitializeGit() error {
	// Get Output path
	outputPath := utils.OutputDir
	// Initialize git repo
	r, err := git.PlainInit(outputPath, false)
	utils.SetGitLocalRepo(r)
	// Check Err
	if err != nil {
		return fmt.Errorf("An Error occurred while initializing repository: %s", err)
	}
	
	// Get working tree
	worktree, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("An error occurred while getting work tree: %v", err)
	} // Debugging: Ensure worktree is not nil
	if worktree == nil {
		return fmt.Errorf("Failed to get worktree: worktree is nil")
	}
	// Set the working tree
	utils.SetGitWorkingTree(worktree)
	fmt.Println("Git working tree initialized successfully at:", outputPath)
	return nil
}

func GitAddFile(filePath string) error {
	// Check if App Git Initilixation is enabled
	if utils.AppGitInit != "yes" {
		return nil
	}
	// Get current worktree
	worktree := utils.GitWorkTree

	// Debugging: Check if worktree is nil
	if worktree == nil {
		return fmt.Errorf("Worktree is nil, ensure InitializeGit() was called successfully")
	}

	fmt.Printf("Attempting to add file to Git: %s\n", filePath)
	relPath, err := filepath.Rel(utils.OutputDir, filePath)
	if err != nil {
		return err
	}
	if relPath == "" {
		return fmt.Errorf("relative path for %s is empyt", filePath)
	}
	fmt.Printf("Adding file to Git (relative path): %s \n", relPath)
	// Add File
	_, err = worktree.Add(relPath)
	// Check for err
	if err != nil {
		return fmt.Errorf("An error occurred while adding file %s to git: %v", filePath, err)
	}
	return nil
}

func GitCommitFiles(message string) error {
	// Check if App Git Initilixation is enabled
	if utils.AppGitInit != "yes" {
		return nil
	}
	// Get Current work tree
	worktree := utils.GitWorkTree
	// Commit file
	_, err := worktree.Commit(message, &git.CommitOptions{})
	// Check for error
	if err != nil {
		return fmt.Errorf("An error occurred while commiting: %v", err)
	}
	return nil
}

func GitPushToRemote() error {
	// Get current local repo
	repo := utils.GitLocalRepo
	// Set Remote Path
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{utils.GitRemoteUrl},
	})
	if err != nil {
		return fmt.Errorf("An error occurred while adding remote config: %v", err)
	}
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: utils.GithubUsername,
			Password: utils.GithubToken,
		},
	})
	if err != nil {
		return fmt.Errorf("An error occurred while pushing online: %v", err)
	}

	return nil
}
