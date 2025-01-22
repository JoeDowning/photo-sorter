package file_manager

import (
	"fmt"
	"os"
)

func GetFiles(path string) (map[string]string, error) {
	entries, err := getDirectoryEntries(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory entries: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() {

		}
	}

	return nil
}

func getDirectoryEntries(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return entries, nil
}
