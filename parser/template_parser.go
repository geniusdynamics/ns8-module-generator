package parser

import (
	"fmt"
	"io/fs"
	"ns8-module-generator/git"
	"os"
	"path/filepath"
	"strings"
)

/*
List Files Read the files in the template folder
*/
func listFiles(folder string) error {
	return filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Print whether it is a file or a directory
		if d.IsDir() {
			fmt.Printf("Directory: %s\n", path)
		} else {
			fmt.Printf("File: %s\n", path)
		}

		return nil
	})
}

// SearchFileAndWriteContent Search for a file in a folder and write content to it
func SearchFileAndWriteContent(folder string, filename string, content string) error {
	var filePath string
	/*
		Search for the file in the folder*
	*/
	err := filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check current entry is a file
		if !d.IsDir() && d.Name() == filename {
			filePath = path
			return filepath.SkipDir
		}
		return nil
	})
	// Check if an error occurred while searching for the file
	if err != nil {
		// if the error is not a file not found error, create the file

		return fmt.Errorf("error while searching for the file: %v", err)
	}
	// Check if the file was found
	if filePath == "" {
		return fmt.Errorf("file not found: %s", filename)
	}
	// Write the content to the file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	// Return nil if everything is successful
	fmt.Printf("Content written to file: %s\n", filePath)

	return nil
}

// SearchFileAndReplaceContent /* SearchFileAndReplaceContent Search for a file in a folder and read content from it and replace placeholders
// replacements is a map of placeholders and their replacements
// map[string]string{"{{placeholder}}": "replacement"}
func SearchFileAndReplaceContent(
	filePath string,
	replacements map[string]string,
) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("File Path does not exist: %s", filePath)
	}
	// Read the content of the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error while reading the file: %v", err)
	}
	// Convert content to string for processing
	content := string(fileContent)
	// Replace placeholders in the content
	for placeholder, replacement := range replacements {
		content = strings.ReplaceAll(content, placeholder, replacement)
	}
	// Write the content bak to the file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	// return nil
	return nil
}

// ReplaceInAllFiles Replace placeholders in all files in a directory
// Such as kickstart which are in all files
func ReplaceInAllFiles(directory string, replacements map[string]string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Process Regular Files
		if !info.IsDir() {
			fmt.Println("The file content of: " + path + " Are being replaced")
			// Call SearchFileAndReplaceContent
			errr := SearchFileAndReplaceContent(path, replacements)
			if errr != nil {
				fmt.Printf("An error occurred: %v", errr)
				return errr
			}
			fmt.Printf("Replaced Content in file: %s\n", path)
		}
		err = git.GitAddFile(path)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
