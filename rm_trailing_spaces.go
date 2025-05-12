package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Starting directory - you can change this to the directory you want to process
	rootDir := "./knowledge"

	// Counter for the number of files processed and modified
	filesProcessed := 0
	filesModified := 0

	// Walk through the directory structure recursively
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process all files (you can add filters here if needed)
		modified, err := removeTrailingSpacesAndBlankLines(path)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", path, err)
			return nil // Continue with other files
		}

		filesProcessed++
		if modified {
			filesModified++
			fmt.Printf("Removed trailing spaces and/or blank lines in: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	fmt.Printf("\nProcessing completed successfully!\n")
	fmt.Printf("Files processed: %d\n", filesProcessed)
	fmt.Printf("Files modified (had trailing spaces or blank lines): %d\n", filesModified)
}

func removeTrailingSpacesAndBlankLines(filePath string) (bool, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("error reading file: %v", err)
	}

	// Convert to string for line processing
	contentStr := string(content)

	// Check if the file is binary (contains null bytes)
	if strings.Contains(contentStr, "\x00") {
		// Skip binary files
		return false, nil
	}

	// Split into lines
	lines := strings.Split(contentStr, "\n")

	// Track if any line was modified
	modified := false

	// Store non-blank lines with trailing spaces removed
	var newLines []string

	// Process each line
	for _, line := range lines {
		// Remove trailing spaces
		trimmedLine := strings.TrimRight(line, " \t")

		// Skip completely blank lines
		if trimmedLine == "" {
			modified = true
			continue
		}

		newLines = append(newLines, trimmedLine)

		if trimmedLine != line {
			modified = true
		}
	}

	// If no lines were modified and no blank lines were removed
	if !modified {
		// Still check for missing trailing newline
		if !strings.HasSuffix(contentStr, "\n") {
			err := ioutil.WriteFile(filePath, append(content, '\n'), 0644)
			if err != nil {
				return false, fmt.Errorf("error appending newline to file: %v", err)
			}
			return true, nil
		}
		return false, nil
	}

	// Join lines back together and add final newline
	newContent := strings.Join(newLines, "\n") + "\n"

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return false, fmt.Errorf("error writing file: %v", err)
	}

	return true, nil
}
