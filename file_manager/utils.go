package file_manager

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var folderFormat = "%s-%s-%s"

func isUsableFileType(fileTypes []string, name string) bool {
	splitName := strings.Split(name, ".")
	if len(splitName) < 2 {
		return false
	}

	for _, fileType := range fileTypes {
		if strings.ToLower(fileType) == strings.ToLower(splitName[1]) {
			return true
		}
	}
	return false
}

func getDirectoryEntries(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return entries, nil
}

func formatFolderName(year, month, day int) string {
	m := strconv.Itoa(month)
	if len(m) == 1 {
		m = "0" + m
	}

	d := strconv.Itoa(day)
	if len(d) == 1 {
		d = "0" + d
	}

	return fmt.Sprintf(folderFormat, strconv.Itoa(year), m, d)
}
