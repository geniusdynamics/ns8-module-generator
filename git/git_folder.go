package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ns8-module-generator/config"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func InitializeGit() error {
	// Get Output path
	outputPath := config.Cfg.OutputDir
	// Initialize git repo

	r, err := git.PlainInit(outputPath, false)
	config.Cfg.GitLocalRepo = r
	// Check Err
	if err != nil {
		return fmt.Errorf("an Error occurred while initializing repository: %s", err)
	}

	headRefs := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main"))

	if err := r.Storer.SetReference(headRefs); err != nil {
		return fmt.Errorf("failed to set HEAD to main: %w", err)
	}
	// Get working tree
	worktree, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("an error occurred while getting work tree: %v", err)
	} // Debugging: Ensure worktree is not nil
	if worktree == nil {
		return fmt.Errorf("failed to get worktree: worktree is nil")
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

	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	}

	if strings.ToLower(config.Cfg.GitAuthMethod) == "ssh" {
		auth, err := sshAuth() // Use the new sshAuth function
		if err != nil {
			return fmt.Errorf("An error occurred while setting up SSH authentication: %v", err)
		}
		pushOptions.Auth = auth
		fmt.Println("Pushing using SSH authentication...")
	} else {
		pushOptions.Auth = &http.BasicAuth{
			Username: config.Cfg.GithubToken,
			Password: config.Cfg.GithubToken,
		}
	}

	err = repo.Push(pushOptions)
	if err != nil {
		return fmt.Errorf("An error occurred while pushing online: %v", err)
	}

	return nil
}

func sshAuth() (ssh.AuthMethod, error) {
	sshAgent, err := ssh.NewSSHAgentAuth("git")
	if err == nil {
		return sshAgent, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home directory: %w", err)
	}

	keyPaths := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_dsa"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
	}

	for _, keyPath := range keyPaths {
		if _, err := os.Stat(keyPath); err == nil {
			publicKeys, err := ssh.NewPublicKeysFromFile("git", keyPath, "")
			if err != nil {
				fmt.Printf("Failed to load SSH key from %s: %v\n", keyPath, err)
				continue // Try next key if loading fails
			}
			return publicKeys, nil
		}
	}

	return nil, fmt.Errorf("no SSH agent found and no suitable SSH key found in common locations")
}
