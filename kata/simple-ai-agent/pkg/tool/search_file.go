package tool

import (
    "fmt"
    "os"
    "path/filepath"
)

// SearchFileAndRead searches for a file by name in the given directory and returns its content.
func SearchFileAndRead(fileName, directory string) (string, error) {
    var filePath string
    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && info.Name() == fileName {
            filePath = path
            return filepath.SkipDir
        }
        return nil
    })
    if err != nil {
        return "", fmt.Errorf("error searching for file: %v", err)
    }
    if filePath == "" {
        return "", fmt.Errorf("file %s not found in %s", fileName, directory)
    }

    content, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("error reading file %s: %v", filePath, err)
    }
    return string(content), nil
}