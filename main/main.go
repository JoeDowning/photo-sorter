package main

import (
	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/logging"

	"go.uber.org/zap"
)

var (
	imageFileTypes               = []string{"jpg", "jpeg", "raw", "cr3", "cr2", "png"}
	sourcePath                   = "/Users/joe.downing/Pictures/Photos/testing-folder"
	destinationPath              = "/Users/joe.downing/Pictures/Photos/testing-folder/sorted"
	nonImageFilesDestinationPath = "/Users/joe.downing/Pictures/Photos/testing-folder/non-image-files/"
)

func main() {
	logger := logging.NewLogger()
	logger = logger.With(
		zap.String("sourcePath", sourcePath),
		zap.String("destinationPath", destinationPath),
		zap.Strings("imageFileTypes", imageFileTypes),
		zap.Bool("includeFiles", true))

	imageFiles, err := file_manager.GetFiles(logger, sourcePath, imageFileTypes, true, image_manager.GetPhoto)
	if err != nil {
		logger.Fatal("failed to get image files", zap.Error(err))
	}

	logger.Debug("got image files", zap.Any("imageFiles", imageFiles))
	sortedFolders := file_manager.SortFilesByDate(imageFiles, image_manager.GetTimestamp)

	logger.Debug("sorted files by date", zap.Any("sortedFolders", sortedFolders))
	err = file_manager.CreateFolderIfNotExists(destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", destinationPath),
			zap.Error(err))
	}

	for folderName, files := range sortedFolders {

		//todo: change so that it is folders of year, month, date

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

	//nonRecognisedFiles, err := file_manager.GetFiles(
	//	logger,
	//	sourcePath,
	//	imageFileTypes,
	//	false,
	//	func(path string) (string, error) { return path, nil },
	//)
	//
	//err = file_manager.CreateFolderIfNotExists(nonImageFilesDestinationPath)
	//if err != nil {
	//	logger.Fatal("failed to create folder in destination path",
	//		zap.String("folderName", nonImageFilesDestinationPath),
	//		zap.Error(err))
	//}
	//
	//for _, filePath := range nonRecognisedFiles {
	//	if strings.Contains(filePath, "DS_Store") {
	//		continue
	//	}
	//	err := file_manager.CopyAndRenameFile(
	//		filePath,
	//		destinationPath+nonImageFilesDestinationPath+filePath)
	//	if err != nil {
	//		logger.Fatal("failed to copy and rename file",
	//			zap.String("destination", destinationPath+nonImageFilesDestinationPath+filePath),
	//			zap.String("file", filePath),
	//			zap.Error(err))
	//	}
	//}
}
