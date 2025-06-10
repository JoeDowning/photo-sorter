package file_manager

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/photos-sorter/pkg/genutils"
)

var folderFormat = "%s-%s-%s"

func getDirectoryEntries(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return entries, nil
}

func formatFolderName(year, month, day int) string {
	return fmt.Sprintf(folderFormat, strconv.Itoa(year),
		genutils.PrefixZeros(2, strconv.Itoa(month)),
		genutils.PrefixZeros(2, strconv.Itoa(day)))
}

func mergeMaps[T any](m1 map[string]T, m2 map[string]T) map[string]T {
	m := make(map[string]T)
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}

// isUsableFileType checks if the file type is in the list of file types, if includeFiles is true
// it will return true if the file type is in the list, if includeFiles is false it will return true
// if the file type is not in the list.
func isUsableFileType(fileTypes []string, name string, includeFiles bool) bool {
	splitName := strings.Split(name, ".")
	if len(splitName) < 2 {
		return false
	}

	var result bool
	for _, fileType := range fileTypes {
		fmt.Printf("Checking file type: %s against %s\n", fileType, splitName[len(splitName)-1])
		if strings.ToLower(fileType) == strings.ToLower(splitName[len(splitName)-1]) && includeFiles {
			result = true
			fmt.Printf("File type %s is usable for file %s\n", fileType, name)
		}
	}
	if !includeFiles && !result {
		result = true
	}
	return result
}
