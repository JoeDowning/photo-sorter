package file_manager

import (
	"fmt"
	"time"
)

func GetFiles[T any](path string, fileTypes []string, fileData func(string) (T, error)) (map[string]T, error) {
	entries, err := getDirectoryEntries(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory entries: %w", err)
	}

	files := make(map[string]T)
	for _, e := range entries {
		if e.IsDir() || !isUsableFileType(fileTypes, e.Name()) {
			continue
		}

		file, err := fileData(path + "/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to get file data: %w", err)
		}
		files[e.Name()] = file
	}

	return files, nil
}

func SortFilesByDate[T any](files map[string]T, getTimeStamp func(T) time.Time) map[string][]T {
	sortedFiles := make(map[string][]T)
	for _, f := range files {
		timestamp := getTimeStamp(f)
		folderName := formatFolderName(timestamp.Year(), int(timestamp.Month()), timestamp.Day())
		if _, ok := sortedFiles[folderName]; !ok {
			sortedFiles[folderName] = make([]T, 0)
		}
		sortedFiles[folderName] = append(sortedFiles[folderName], f)
	}

	return sortedFiles
}
