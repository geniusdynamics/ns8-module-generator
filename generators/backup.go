package generators

import (
	"fmt"
	"ns8-module-generator/git"
	"os"
)

func AddToBackup(filePath, backupContent string) error {
	err := WriteToFile(filePath, backupContent, "feat(backup): added backup files needed")
	return err
}

func WriteBackUpScript(filePath, backupContent string) error {
	err := WriteToFile(filePath, backupContent, "feat: Added Backup COntent")
	return err
}

func WriteToFile(filePath, fileContent, commitMessage string) error {
	content, err := os.ReadFile(filePath)
	// If error occurs Close
	if err != nil {
		return fmt.Errorf("failed to read the file: %v", err)
	}
	fmt.Println("existing file contents: ", string(content))

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v ", err)
	}
	// Close file later
	defer file.Close()
	if _, err := file.WriteString(fileContent + "\n"); err != nil {
		return fmt.Errorf("failed to add json dump in %s;", filePath)
	}
	fmt.Print(fileContent)
	err = git.GitAddFile(filePath)
	if err != nil {
		return err
	}
	err = git.GitCommitFiles(commitMessage)
	if err != nil {
		return err
	}
	return nil
}
