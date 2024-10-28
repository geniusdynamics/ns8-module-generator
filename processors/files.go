package processors

import (
	"archive/zip"
	"fmt"
	"io"
	"ns8-module-generator/utils"
	"os"
	"path/filepath"
)

// CopyDirectory copies the template directory to the output directory with a progress indicator.
func CopyDirectory() error {
	// Create the output directory
	err := os.MkdirAll(utils.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("error while creating the output directory: %v", err)
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
