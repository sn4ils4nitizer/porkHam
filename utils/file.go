package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// sets directory
const pagesDir = "pages"

// Read file from disk
func ReadFile(name string) (string, error) {
	filePath := filepath.Join(pagesDir, name)

	absPath, _ := filepath.Abs(filePath)
	log.Println("file.go: Attempting to read file: ", absPath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile, save HTML content to a file
func WriteFile(fullpath string, content string) error {
	dir := filepath.Dir(fullpath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return ioutil.WriteFile(fullpath, []byte(content), 0644)
}

// DeleteFile
func DeleteFile(name string) error {
	filePath := pagesDir + name + ".html"
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("faile to delete file: %s: %v", filePath, err)
	}
	return nil
}
