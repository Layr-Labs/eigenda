package client

import (
	"fmt"
	"os"
	"strings"
)

// ResolveTildeInPath resolves the tilde (~) in the given path to the user's home directory.
func ResolveTildeInPath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return strings.Replace(path, "~", homeDir, 1), nil
}
