package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Starting directory
	rootDir := "knowledge"

	// Walk through the directory structure recursively
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has .yaml or .txt extension
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".txt") {
			processFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	fmt.Println("Processing completed successfully!")
}

func processFile(filePath string) {
	fmt.Printf("Processing file: %s\n", filePath)

	// Extract directory and filename
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	// Handle .yaml files
	if strings.HasSuffix(filename, "_qna.yaml") {
		// Extract base name for new directory
		baseName := strings.TrimSuffix(filename, "_qna.yaml")
		newDirPath := filepath.Join(dir, baseName)

		// Create new directory
		err := os.MkdirAll(newDirPath, 0755)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", newDirPath, err)
			return
		}

		// New file path
		newFilePath := filepath.Join(newDirPath, "qna.yaml")

		// Read content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return
		}

		// Edit content: version 1 -> version 3 and commit 1e7334d -> 3e10cd8
		contentStr := string(content)
		contentStr = strings.Replace(contentStr, "version: 1", "version: 3", 1)
		contentStr = strings.Replace(contentStr, "data_md/", "", 1)
		// contentStr = strings.Replace(contentStr, "commit: 1e7334d", "commit: 3e10cd8", 1)

		// Write to new file
		err = ioutil.WriteFile(newFilePath, []byte(contentStr), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", newFilePath, err)
			return
		}

		// Remove original file
		err = os.Remove(filePath)
		if err != nil {
			fmt.Printf("Error removing original file %s: %v\n", filePath, err)
		}

		fmt.Printf("Processed YAML file: %s -> %s\n", filePath, newFilePath)
	} else if strings.HasSuffix(filename, ".txt") {
		// For .txt files
		// Extract base name (assuming format might be something.txt)
		baseName := strings.TrimSuffix(filename, ".txt")
		newDirPath := filepath.Join(dir, baseName)

		// Check if directory exists, if not create it
		if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
			err := os.MkdirAll(newDirPath, 0755)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", newDirPath, err)
				return
			}
		}

		// New file path
		newFilePath := filepath.Join(newDirPath, "attribution.txt")

		// Read content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return
		}

		// Write to new file
		err = ioutil.WriteFile(newFilePath, content, 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", newFilePath, err)
			return
		}

		// Remove original file
		err = os.Remove(filePath)
		if err != nil {
			fmt.Printf("Error removing original file %s: %v\n", filePath, err)
		}

		fmt.Printf("Processed TXT file: %s -> %s\n", filePath, newFilePath)
	}
}
