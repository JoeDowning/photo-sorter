package main

import (
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/config"
	"github.com/photos-sorter/pkg/logging"
	"github.com/photos-sorter/sorting"
)

const (
	videoMode = "videos"
	imageMode = "images"

	moveFileMode = "move"
	copyFileMode = "copy"

	mode     = "images"
	fileMode = "copy"
)

//todo look at uploading to google photos
//todo test the moving of videos
//todo double check it won't try copying a file that is already there (log this)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to get config", zap.Error(err))
	}

	logger := logging.NewLogger(cfg.LogLevel)
	logger.Info("Started photos sorter",
		zap.String("sourcePath", cfg.SourcePath),
		zap.String("destinationPath", cfg.DestinationPath),
		zap.String("mode", mode),
		zap.String("fileMode", fileMode),
		zap.Bool("includeFiles", true))

	startTime := time.Now()

	var moveFileFunc func(*zap.Logger, string, string) error
	switch cfg.FileMode {
	case moveFileMode:
		moveFileFunc = file_manager.MoveAndRenameFile
	case copyFileMode:
		moveFileFunc = file_manager.CopyAndRenameFile
	default:
		logger.Fatal("invalid file mode selected", zap.String("fileMode", fileMode))
	}

	switch cfg.FileType {
	case imageMode:
		err = sorting.SortImages(logger, cfg, moveFileFunc)
	case videoMode:
		err = sorting.SortVideos(logger, cfg)
	default:
		logger.Fatal("invalid mode selected", zap.String("mode", cfg.FileMode))
	}
	if err != nil {
		logger.Fatal("failed to sort files", zap.Error(err))
	}

	logger.Info("Finished photos sorter",
		zap.Duration("runTime", time.Since(startTime)),
		zap.Int("filesMoved", file_manager.ReturnFilesCount()),
		zap.Int("entriesChecked", file_manager.ReturnEntriesCheckedCount()))
}

func usingSortedFolders(logger *zap.Logger, cfg config.Config, imageFiles map[string]image_manager.ImageData) {
	sortedFolders := file_manager.SortFilesByDate(imageFiles, image_manager.GetTimestamp)

	logger.Debug("sorted files by date", zap.Any("sortedFolders", sortedFolders))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.DestinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.DestinationPath),
			zap.Error(err))
	}

	for folderName, files := range sortedFolders {
		err := file_manager.CreateFolderIfNotExists(logger, cfg.DestinationPath+"/"+folderName)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", folderName),
				zap.Error(err))
		}

		for _, file := range files {
			err := file_manager.MoveAndRenameFile(
				logger,
				file.GetFilePath(),
				cfg.DestinationPath+"/"+folderName+"/"+file.GetFileName())
			if err != nil {
				logger.Fatal("failed to copy and rename file",
					zap.String("destination", cfg.DestinationPath+"/"+folderName+"/"+file.GetFileName()),
					zap.String("file", file.GetFileName()),
					zap.Error(err))
			}
		}
	}
}

//func nonRecognisedFileSorter(logger *zap.Logger) {
//	nonRecognisedFiles, err := file_manager.GetFilesSingleFolder(
//		logger,
//		testSourcePath,
//		image_manager.GetImageTypes(),
//		false,
//		func(path string) (string, error) { return path, nil },
//	)
//
//	err = file_manager.CreateFolderIfNotExists(logger, testNonImageFilesDestinationPath)
//	if err != nil {
//		logger.Fatal("failed to create folder in destination path",
//			zap.String("folderName", testNonImageFilesDestinationPath),
//			zap.Error(err))
//	}
//
//	for _, filePath := range nonRecognisedFiles {
//		if strings.Contains(filePath, "DS_Store") {
//			continue
//		}
//		err := file_manager.MoveAndRenameFile(
//			logger,
//			filePath,
//			testDestinationPath+testNonImageFilesDestinationPath+filePath)
//		if err != nil {
//			logger.Fatal("failed to copy and rename file",
//				zap.String("destination", testDestinationPath+testNonImageFilesDestinationPath+filePath),
//				zap.String("file", filePath),
//				zap.Error(err))
//		}
//	}
//}
