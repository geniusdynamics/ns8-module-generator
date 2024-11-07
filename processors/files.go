package processors

import (
	"archive/zip"
	"fmt"
	"io"
	"ns8-module-generator/git"
	"ns8-module-generator/utils"
	"os"
	"path/filepath"
	"strings"
)

// CopyDirectory copies the template directory to the output directory with a progress indicator.
func CopyDirectory() error {
	// Create the output directory
	err := os.MkdirAll(utils.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("error while creating the output directory: %v", err)
	}
	if utils.AppGitInit == "yes" {
		// Initialize git repo
		err = git.InitializeGit()
		if err != nil {
			fmt.Printf("An error occurred while initializing git: %s", err)
		}
	}

	// Count total number of files and directories to copy
	var totalItems int
	err = filepath.WalkDir(utils.TemplateDir, func(_ string, info os.DirEntry, _ error) error {
		if !info.IsDir() {
			totalItems++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error counting items in the source directory: %v", err)
	}

	// Initialize progress variables
	var copiedItems int

	// Function to print progress
	printProgress := func() {
		progress := float64(copiedItems) / float64(totalItems) * 100
		fmt.Printf("\rProgress: %.2f%% (%d/%d)", progress, copiedItems, totalItems)
	}

	// Walk through the source directory
	err = filepath.WalkDir(
		utils.TemplateDir,
		func(srcPath string, info os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking through the source directory: %v", err)
			}

			// Create the destination path
			relPath, _ := filepath.Rel(utils.TemplateDir, srcPath)
			dstPath := filepath.Join(utils.OutputDir, relPath)

			if info.IsDir() {
				// Create directory in the destination
				err := os.MkdirAll(dstPath, 0755)
				if err != nil {
					return fmt.Errorf("failed to create directory %s: %v", dstPath, err)
				}
			} else {
				// Copy file
				srcFile, err := os.Open(srcPath)
				if err != nil {
					return fmt.Errorf("failed to open source file %s: %v", srcPath, err)
				}
				defer srcFile.Close()

				dstFile, err := os.Create(dstPath)
				if err != nil {
					return fmt.Errorf("failed to create destination file %s: %v", dstPath, err)
				}
				defer dstFile.Close()

				_, err = io.Copy(dstFile, srcFile)
				if err != nil {
					return fmt.Errorf("failed to copy file content from %s to %s: %v", srcPath, dstPath, err)
				}
				// Git add File
				err = git.GitAddFile(dstPath)
				if err != nil {
					return fmt.Errorf("failed to add file to git worktree: %s", err)
				}
			}

			// Update progress
			copiedItems++
			printProgress()
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("error copying directory: %v", err)
	}

	// Complete progress output
	fmt.Printf("\nCopy complete! %d items copied.\n", copiedItems)
	return nil
}

// ZipOutput Zip the output directory to a zip file
func ZipOutput(name string) error {
	// Create the zip file
	fmt.Println("Creating zip file...")
	zipFile, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("error creating zip file: %v", err)
	}
	defer zipFile.Close()
	println("Zip file created")
	zipWrite := zip.NewWriter(zipFile)
	defer zipWrite.Close()

	// Walk Through the output directory
	err = filepath.WalkDir(utils.OutputDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through the output directory: %v", err)
		}
		// Create the file in the zip
		relPath, _ := filepath.Rel(utils.OutputDir, path)
		if info.IsDir() {
			relPath += "/"
		}
		zipFile, err := zipWrite.Create(relPath)
		if err != nil {
			return fmt.Errorf("error creating file in zip: %v", err)
		}
		if info.IsDir() {
			return nil
		}
		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()
		// Copy the file to the zip
		_, err = io.Copy(zipFile, file)
		if err != nil {
			return fmt.Errorf("error copying file to zip: %v", err)
		}
		return nil
	})

	return nil
}

// CleanOutputDirectory Clean the output directory
func CleanOutputDirectory() error {
	err := os.RemoveAll(utils.OutputDir)
	if err != nil {
		return fmt.Errorf("error cleaning output directory: %v", err)
	}
	return nil
}

func CreateModuleFromTemplateDirectory(name string) error {
	err := CopyDirectory()
	if err != nil {
		return fmt.Errorf("error copying directory: %v", err)
	}
	err = ZipOutput(name + ".zip")
	if err != nil {
		return fmt.Errorf("error zipping output: %v", err)
	}
	err = CleanOutputDirectory()
	if err != nil {
		return fmt.Errorf("error cleaning output directory: %v", err)
	}
	return nil
}

func UnzipFiles(destPath, zipPath string) {
	// Ensure the destination directory exists; create if it doesn't
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		fmt.Printf("Failed to create destination directory: %v\n", err)
		return
	}
	// Open the zip archive
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("An error occurred while unzipping files: %v\n", err)
		return
	}
	defer archive.Close()

	// Loop through files in the archive
	for _, f := range archive.File {
		// Strip out any root directory from each file's name to avoid nested directories
		filePath := filepath.Join(destPath, f.Name)

		// Ensure file paths are valid and inside the destination directory
		if !strings.HasPrefix(filePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			fmt.Println("Invalid file path")
			return
		}

		// Create directories as needed
		if f.FileInfo().IsDir() {
			fmt.Println("Creating directory...")
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", filePath, err)
				return
			}
			continue
		}

		// Create the file
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", filepath.Dir(filePath), err)
			return
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			return
		}
		defer dstFile.Close()

		// Open the file within the archive
		fileInArchive, err := f.Open()
		if err != nil {
			fmt.Printf("Error opening file in archive %s: %v\n", f.Name, err)
			return
		}
		defer fileInArchive.Close()

		// Copy contents to the destination file
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			fmt.Printf("Error copying contents to file %s: %v\n", filePath, err)
			return
		}
	}
	fmt.Println("Unzipping complete!")
	err = os.Remove(zipPath)
	if err != nil {
		fmt.Print("An error occurred while deleting zip")
	}
}
