package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/logging"

	"go.uber.org/zap"
)

var (
	imageFileTypes               = []string{"jpg", "jpeg", "raw", "cr3", "cr2", "png"}
	rawFileTypes                 = []string{"raw", "cr3", "cr2"}
	editedFilesContainText       = []string{"tz", "ps", "dxo", "e"}
	sourcePath                   = "/Users/joe.downing/Pictures/Photos/testing-folder/test-images"
	destinationPath              = "/Users/joe.downing/Pictures/Photos/testing-folder/sorted"
	nonImageFilesDestinationPath = "/Users/joe.downing/Pictures/Photos/testing-folder/non-image-files/"
)

//todo: add a function to return the camera model
//todo: ignore anything that isn't a canon or panasonic camera?
//todo: use the raw/edited function

func main() {
	logger := logging.NewLogger()
	logger.Info("Started photos sorter",
		zap.String("sourcePath", sourcePath),
		zap.String("destinationPath", destinationPath),
		zap.Strings("imageFileTypes", imageFileTypes),
		zap.Bool("includeFiles", true))
	startTime := time.Now()

	//imageFiles, err := file_manager.GetFilesSingleFolder(logger, sourcePath, imageFileTypes, true, image_manager.GetPhoto)
	//if err != nil {
	//	logger.Fatal("failed to get image files", zap.Error(err))
	//}

	imageFiles, err := file_manager.GetFilesAllDepths(logger, sourcePath, imageFileTypes, true, image_manager.GetPhoto)
	if err != nil {
		logger.Fatal("failed to get image files", zap.Error(err))
	}

	logger.Info("Got image files", zap.Int("count", len(imageFiles)))

	// enable for sorting into folder structure of "year/month/day/<file>"
	usingFilesWithPath(logger, imageFiles)

	logger.Info("Finished photos sorter", zap.Duration("runTime", time.Since(startTime)))
	// enable for sorting into folder struct of "year-month-day/<file>"
	//usingSortedFolders(logger, imageFiles)

	// enable for moving non-recognised files to another folder
	//nonRecognisedFileSorter(logger)
}

func usingSortedFolders(logger *zap.Logger, imageFiles map[string]image_manager.ImageData) {
	sortedFolders := file_manager.SortFilesByDate(imageFiles, image_manager.GetTimestamp)

	logger.Debug("sorted files by date", zap.Any("sortedFolders", sortedFolders))
	err := file_manager.CreateFolderIfNotExists(destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", destinationPath),
			zap.Error(err))
	}

	for folderName, files := range sortedFolders {
		err := file_manager.CreateFolderIfNotExists(destinationPath + "/" + folderName)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", folderName),
				zap.Error(err))
		}

		for _, file := range files {
			err := file_manager.CopyAndRenameFile(
				file.GetFilePath(),
				destinationPath+"/"+folderName+"/"+file.GetFileName())
			if err != nil {
				logger.Fatal("failed to copy and rename file",
					zap.String("destination", destinationPath+"/"+folderName+"/"+file.GetFileName()),
					zap.String("file", file.GetFileName()),
					zap.Error(err))
			}
		}
	}
}

func usingFilesWithPath(logger *zap.Logger, imageFiles map[string]image_manager.ImageData) {
	logger.Info("Sorting files using source paths", zap.String("destinationPath", destinationPath))
	err := file_manager.CreateFolderIfNotExists(destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", destinationPath),
			zap.Error(err))
	}

	filesWithPath := file_manager.AddFolderPathToFile(
		imageFiles,
		func(file image_manager.ImageData) image_manager.ImageData {
			timestamp := image_manager.GetTimestamp(file)
			year := strconv.Itoa(timestamp.Year())
			month := strconv.Itoa(int(timestamp.Month()))
			if len(month) == 1 {
				month = "0" + month
			}
			day := strconv.Itoa(timestamp.Day())
			if len(day) == 1 {
				day = "0" + day
			}
			file.DestPath = year + "/" + month + "/" + day + "/" + file.GetFileName()
			return file
		})

	for _, file := range filesWithPath {
		logger.Debug("copying file",
			zap.String("destination", destinationPath+"/"+file.DestPath),
			zap.String("file", file.GetFileName()))

		err := file_manager.CreatePathFoldersIfDoesntExists(destinationPath, file.DestPath)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", file.GetFilePath()),
				zap.Error(err))
		}

		err = file_manager.CopyAndRenameFile(
			file.GetFilePath(),
			destinationPath+"/"+file.DestPath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", destinationPath+"/"+file.DestPath),
				zap.String("file", file.GetFileName()),
				zap.Error(err))
		}
	}
}

func nonRecognisedFileSorter(logger *zap.Logger) {
	nonRecognisedFiles, err := file_manager.GetFilesSingleFolder(
		logger,
		sourcePath,
		imageFileTypes,
		false,
		func(path string) (string, error) { return path, nil },
	)

	err = file_manager.CreateFolderIfNotExists(nonImageFilesDestinationPath)
	if err != nil {
		logger.Fatal("failed to create folder in destination path",
			zap.String("folderName", nonImageFilesDestinationPath),
			zap.Error(err))
	}

	for _, filePath := range nonRecognisedFiles {
		if strings.Contains(filePath, "DS_Store") {
			continue
		}
		err := file_manager.CopyAndRenameFile(
			filePath,
			destinationPath+nonImageFilesDestinationPath+filePath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", destinationPath+nonImageFilesDestinationPath+filePath),
				zap.String("file", filePath),
				zap.Error(err))
		}
	}
}

func editedOrRawPath(fullFileName string) string {
	splitName := strings.Split(fullFileName, ".")
	if len(splitName) < 2 {
		return ""
	}

	fileType := strings.ToLower(splitName[1])
	if file_manager.InArray(rawFileTypes, fileType) {
		return "raw"
	}

	fileName := strings.ToLower(splitName[0])
	if file_manager.ContainsFromArray(editedFilesContainText, fileName) {
		return "edited"
	}

	return "raw"
}
