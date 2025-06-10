package file_manager

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	entriesCheckedCount int
	movedFileCount      int
	FilesToMoveCount    int
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
		if e.IsDir() || !isUsableFileType(fileTypes, e.Name(), includeFiles) {
			logger.Debug("skipping file",
				zap.String("name", e.Name()),
				zap.Bool("isDir", e.IsDir()),
				zap.Bool("isUsableFileType", isUsableFileType(fileTypes, e.Name(), includeFiles)),
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
	includeFiles bool, fileData func(*zap.Logger, string) (T, error)) (map[string]T, error,
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
		} else if isUsableFileType(fileTypes, e.Name(), includeFiles) {
			logger.Debug("getting file data", zap.String("name", e.Name()))
			file, err := fileData(logger, path+"/"+e.Name())
			if err != nil {
				logger.Error("failed to get file data",
					zap.String("name", e.Name()))
				continue
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

func AddFolderPathToFile[T any](logger *zap.Logger, files map[string]T, addFolderPath func(*zap.Logger, T) T) map[string]T {
	filesWithPath := make(map[string]T)
	for name, file := range files {
		fileWithPath := addFolderPath(logger, file)
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

func MoveAndRenameFile(logger *zap.Logger, src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		logger.Debug("Destination file already exists", zap.String("destination", dst))
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check destination file: %w", err)
	}

	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	movedFileCount++
	logger.Info(fmt.Sprintf("[ %d / %d ] files moved", movedFileCount, FilesToMoveCount))
	return nil
}

func CopyAndRenameFile(logger *zap.Logger, src, dst string) error {
	err := copyFile(logger, src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	movedFileCount++
	logger.Info(fmt.Sprintf("[ %d / %d ] files copied", movedFileCount, FilesToMoveCount))
	return nil
}

func copyFile(logger *zap.Logger, src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		logger.Debug("Destination file already exists", zap.String("destination", dst))
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check destination file: %w", err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func ReturnFilesCount() int {
	return movedFileCount
}

func ReturnEntriesCheckedCount() int {
	return entriesCheckedCount
}
