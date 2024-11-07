package generators

import (
	"fmt"
	"ns8-module-generator/git"
	"os"
)

func AddToBackup(filePath, backupContent string) error {
	// Check if file Exists
	content, err := os.ReadFile(filePath)
	// If error occurs Close
	if err != nil {
		return fmt.Errorf("Failed to read the file: %v", err)
	}
	fmt.Println("Existing File Contents: ", string(content))

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write to file: %v ", err)
	}
	// Close file later
	defer file.Close()
	if _, err := file.WriteString(backupContent + "\n"); err != nil {
		return fmt.Errorf("Failed to add JSON DUMP in %s;", filePath)
	}
	err = git.GitAddFile(filePath)
	if err != nil {
		return err
	}
	err = git.GitCommitFiles("feat(backup): added backup files needed")
	if err != nil {
		return err
	}
	return nil
}
