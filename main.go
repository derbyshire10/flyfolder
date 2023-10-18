package main

import (
	"fmt"
	"flag"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Define flags and set default values
	selectedFolderPtr := flag.String("folder", "", "Folder selected to duplicate")

	// Parse the command-line arguments
	flag.Parse()

	// Check if required flags are provided 
	if selectedFolderPtr == nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get the current directory path
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	// Open and read the current directory
	dir, err := os.Open(currentDir)
	if err != nil {
		fmt.Println("Error opening current directory:", err)
		return
	}
	defer dir.Close()

	// Read the directory entries
	entries, err := dir.Readdir(0)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Iterate through the entries and check for a match with the selected directory
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == *selectedFolderPtr {
				destDir := entry.Name() + "1"
				err = copyDir(entry.Name(), destDir)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Directory copied successfully.")
				}
			}
		}
	}
}

func copyDir(srcDir, destDir string) error {
	// Create the destination directory
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	// Open the source directory
	src, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the source directory contents
	entries, err := src.Readdir(-1)
	if err != nil {
		return err
	}

	// Copy each file and subdirectory to the destination
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}