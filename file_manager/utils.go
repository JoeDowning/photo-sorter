package file_manager

import (
	"fmt"
	"os"
	"strconv"

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
