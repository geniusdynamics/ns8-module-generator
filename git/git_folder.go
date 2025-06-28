package git

import (
	"fmt"
	"ns8-module-generator/config"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func InitializeGit() error {
	// Get Output path
	outputPath := config.Cfg.OutputDir
	// Initialize git repo
	r, err := git.PlainInit(outputPath, false)
	config.Cfg.GitLocalRepo = r
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
	config.Cfg.GitWorkTree = worktree
	fmt.Println("Git working tree initialized successfully at:", outputPath)
	return nil
}

func GitAddFile(filePath string) error {
	// Check if App Git Initilixation is enabled
	if !config.Cfg.AppGitInit {
		return nil
	}
	// Get current worktree
	worktree := config.Cfg.GitWorkTree

	// Debugging: Check if worktree is nil
	if worktree == nil {
		return fmt.Errorf("Worktree is nil, ensure InitializeGit() was called successfully")
	}

	fmt.Printf("Attempting to add file to Git: %s\n", filePath)
	relPath, err := filepath.Rel(config.Cfg.OutputDir, filePath)
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
	if !config.Cfg.AppGitInit {
		return nil
	}
	// Get Current work tree
	worktree := config.Cfg.GitWorkTree
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
	repo := config.Cfg.GitLocalRepo
	// Set Remote Path
	_, err := repo.CreateRemote(&gitConfig.RemoteConfig{
		Name: "origin",
		URLs: []string{config.Cfg.GitRemoteUrl},
	})
	if err != nil {
		return fmt.Errorf("An error occurred while adding remote config: %v", err)
	}
	fmt.Print("Your github username: " + config.Cfg.GithubUsername + "\n")
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: config.Cfg.GithubToken,
			Password: config.Cfg.GithubToken,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("An error occurred while pushing online: %v", err)
	}

	return nil
}
