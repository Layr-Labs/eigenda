package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// findGoMod searches for the go.mod file in the current directory and all parent directories.
func findGoMod(path string) (string, error) {
	// Check if go.mod exists in the current directory
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return path, nil
	}

	// Get the parent directory
	parent := filepath.Dir(path)

	// If the parent directory is the same as the current one, we've reached the root of the file system
	if parent == path {
		return "", fmt.Errorf("no go.mod found")
	}

	// Recursively look in the parent directory
	return findGoMod(parent)
}

func MustGetModuleRootPath() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %s\n", err)
		os.Exit(1)
	}

	// Find the directory with go.mod
	rootDir, err := findGoMod(cwd)
	if err != nil {
		panic(fmt.Errorf("Error finding go.mod: %s\n", err))
	}

	return rootDir
}
