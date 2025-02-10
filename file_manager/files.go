package file_manager

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/photos-sorter/sorting"
)

var (
	entriesCheckedCount int
	movedFileCount      int
)

func GetFilesSingleFolder[T any](logger *zap.Logger, path string, fileTypes []string,
	includeFiles bool, fileData func(string) (T, error)) (map[string]T, error,
) {
	entries, err := getDirectoryEntries(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory entries: %w", err)
	}

	logger.Debug("got directory entries",
		zap.String("path", path),
		zap.Any("entries", entries))
	files := make(map[string]T)
	for _, e := range entries {
		entriesCheckedCount++
		logger.Info(fmt.Sprintf("%d entries checked", entriesCheckedCount))
		if e.IsDir() || !sorting.IsUsableFileType(fileTypes, e.Name(), includeFiles) {
			logger.Debug("skipping file",
				zap.String("name", e.Name()),
				zap.Bool("isDir", e.IsDir()),
				zap.Bool("isUsableFileType", sorting.IsUsableFileType(fileTypes, e.Name(), includeFiles)),
				zap.Bool("includeFiles", includeFiles))
			continue
		}

		logger.Debug("getting file data",
			zap.String("name", e.Name()))
		file, err := fileData(path + "/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to get file data: %w", err)
		}
		logger.Debug("got file data",
			zap.String("name", e.Name()),
			zap.Any("file", file))
		files[e.Name()] = file
	}

	return files, nil
}

func GetFilesAllDepths[T any](logger *zap.Logger, path string, fileTypes []string,
	includeFiles bool, fileData func(string) (T, error)) (map[string]T, error,
) {
	entries, err := getDirectoryEntries(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory entries: %w", err)
	}
	logger.Debug("got directory entries", zap.String("path", path), zap.Any("entries", entries))

	var directoryTotal, fileTotal int
	files := make(map[string]T)
	for _, e := range entries {
		entriesCheckedCount++
		logger.Info(fmt.Sprintf("%d entries checked", entriesCheckedCount))

		if e.IsDir() {
			logger.Debug("getting files from subfolder", zap.String("name", e.Name()))
			subFiles, err := GetFilesAllDepths(logger, path+"/"+e.Name(), fileTypes, includeFiles, fileData)
			if err != nil {
				return nil, fmt.Errorf("failed to get files from subfolder: %w", err)
			}

			files = mergeMaps(files, subFiles)
			directoryTotal++
		} else if sorting.IsUsableFileType(fileTypes, e.Name(), includeFiles) {
			logger.Debug("getting file data", zap.String("name", e.Name()))
			file, err := fileData(path + "/" + e.Name())
			if err != nil {
				logger.Error("failed to get file data",
					zap.String("name", e.Name()))
				return nil, fmt.Errorf("failed to get file data: %w", err)
			}
			logger.Debug("got file data", zap.String("name", e.Name()), zap.Any("file", file))

			files[e.Name()] = file
			fileTotal++
		}
	}

	logger.Debug("got files from folder",
		zap.String("path", path),
		zap.Int("directoryTotal", directoryTotal),
		zap.Int("fileTotal", fileTotal),
		zap.Int("totalEntries", directoryTotal+fileTotal))
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

func AddFolderPathToFile[T any](files map[string]T, addFolderPath func(T) T) map[string]T {
	filesWithPath := make(map[string]T)
	for name, file := range files {
		fileWithPath := addFolderPath(file)
		filesWithPath[name] = fileWithPath
	}
	return filesWithPath
}

func CreateFolderIfNotExists(logger *zap.Logger, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Debug("Creating folder", zap.String("folderPath", path))
		err := os.Mkdir(path, 0750)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// CreatePathFoldersIfDoesntExists creates folders for the given path if they don't exist
// foldersPath: the path to the folder where the folders will be created
func CreatePathFoldersIfDoesntExists(logger *zap.Logger, foldersPath, path string) error {
	folders := strings.Split(path, "/")
	for i, folder := range folders {
		if i+1 == len(folders) {
			break
		}
		foldersPath += "/" + folder
		if _, err := os.Stat(foldersPath); os.IsNotExist(err) {
			logger.Debug("Creating folder", zap.String("folderPath", foldersPath))
			err := os.Mkdir(foldersPath, 0750)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}
	}

	return nil
}

func CopyAndRenameFile(logger *zap.Logger, src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	movedFileCount++
	logger.Info(fmt.Sprintf("%d files moved", movedFileCount))
	return nil
}

func ReturnFilesCount() int {
	return movedFileCount
}

func ReturnEntriesCheckedCount() int {
	return entriesCheckedCount
}
