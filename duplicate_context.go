package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Starting directory
	rootDir := "knowledge"

	// Files processed counter
	totalFiles := 0
	modifiedFiles := 0

	// Walk through the directory structure recursively
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process qna.yaml files
		if !strings.HasSuffix(path, "qna.yaml") {
			return nil
		}

		totalFiles++

		// Process the file
		modified, err := processQnaFile(path)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", path, err)
			return nil // Continue with other files
		}

		if modified {
			modifiedFiles++
			fmt.Printf("Modified: %s\n", path)
		} else {
			fmt.Printf("Skipped (already has 5+ contexts): %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	fmt.Printf("\nProcessing completed successfully!\n")
	fmt.Printf("Total qna.yaml files processed: %d\n", totalFiles)
	fmt.Printf("Files modified: %d\n", modifiedFiles)
}

func processQnaFile(filePath string) (bool, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("error reading file: %v", err)
	}

	// Convert to string for processing
	contentStr := string(content)

	// Get all the context sections
	re := regexp.MustCompile(`(?s)(  - context: \|.+?)(?:  - context: \||document_outline: \|)`)
	matches := re.FindAllStringSubmatch(contentStr, -1)

	// Count the context sections
	contextCount := len(matches)
	if contextCount >= 5 {
		return false, nil // No need to modify
	}

	// If no context sections found, return error
	if contextCount == 0 {
		return false, fmt.Errorf("no context sections found in file")
	}

	// Get the last complete context section
	lastContextSection := matches[len(matches)-1][1]

	// Find where to insert the duplicated contexts
	outlineIndex := strings.Index(contentStr, "document_outline: |")
	if outlineIndex == -1 {
		return false, fmt.Errorf("document_outline not found in file")
	}

	// Find the position before document_outline where we will insert new contexts
	// We need to find the last newline before document_outline
	insertPos := strings.LastIndex(contentStr[:outlineIndex], "\n")
	if insertPos == -1 {
		insertPos = 0
	} else {
		insertPos++ // Move past the newline
	}

	// Create the new content
	var newContent strings.Builder
	newContent.WriteString(contentStr[:insertPos])

	// Add duplicated contexts
	for i := contextCount; i < 5; i++ {
		newContent.WriteString(lastContextSection)
	}

	// Add the rest of the file
	newContent.WriteString(contentStr[insertPos:])

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent.String()), 0644)
	if err != nil {
		return false, fmt.Errorf("error writing file: %v", err)
	}

	return true, nil
}
