package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReplaceTextInFile(filePath, oldText, newText string) error {
	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Replace the text within the content
	newContent := strings.ReplaceAll(string(content), oldText, newText)

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceTextInFolder(rootFolder, oldText, newText string) error {
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".git") {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Replace text in the file
		err = ReplaceTextInFile(path, oldText, newText)
		if err != nil {
			fmt.Printf("Error replacing text in file %s: %s\n", path, err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
