package main

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/logging"
	"github.com/photos-sorter/sorting"
	"github.com/photos-sorter/video_manager"
)

const (
	testSourcePath                   = "/Users/joe.downing/Pictures/Photos/testing-folder/test-images"
	testZipSourcePath                = "/Users/joe.downing/Pictures/Photos/testing-folder/test-zips"
	testDestinationPath              = "/Users/joe.downing/Pictures/Photos/testing-folder/sorted"
	testNonImageFilesDestinationPath = "/Users/joe.downing/Pictures/Photos/testing-folder/non-image-files/"

	seagateSourcePath      = "/Volumes/Seagate/takeouts"
	seagateDestinationPath = "/Volumes/Seagate/sorted"
)

const (
	videoMode = "videos"
	imageMode = "images"

	moveFileMode = "move"
	copyFileMode = "copy"

	mode     = "images"
	fileMode = "copy"
)

var (
	testConfig = config{
		sourcePath:      testSourcePath,
		destinationPath: testDestinationPath,
	}
	testZipConfig = config{
		sourcePath:      testZipSourcePath,
		destinationPath: testDestinationPath,
	}
	seagateConfig = config{
		sourcePath:      seagateSourcePath,
		destinationPath: seagateDestinationPath,
	}
)

type config struct {
	sourcePath      string
	destinationPath string
}

//todo look at uploading to google photos
//todo test the moving of videos

func main() {
	cfg := testZipConfig
	logger := logging.NewLogger()

	logger.Info("Started photos sorter",
		zap.String("sourcePath", cfg.sourcePath),
		zap.String("destinationPath", cfg.destinationPath),
		zap.String("mode", mode),
		zap.String("fileMode", fileMode),
		zap.Bool("includeFiles", true))
	startTime := time.Now()

	var moveFileFunc func(*zap.Logger, string, string) error
	switch fileMode {
	case moveFileMode:
		moveFileFunc = file_manager.MoveAndRenameFile
	case copyFileMode:
		moveFileFunc = file_manager.CopyAndRenameFile
	default:
		logger.Fatal("invalid file mode selected", zap.String("fileMode", fileMode))
	}

	var err error
	switch mode {
	case imageMode:
		err = sortImages(logger, cfg, moveFileFunc)
	case videoMode:
		err = sortVideos(logger, cfg)
	default:
		logger.Fatal("invalid mode selected", zap.String("mode", mode))
	}
	if err != nil {
		logger.Fatal("failed to sort files", zap.Error(err))
	}

	logger.Info("Finished photos sorter",
		zap.Duration("runTime", time.Since(startTime)),
		zap.Int("filesMoved", file_manager.ReturnFilesCount()),
		zap.Int("entriesChecked", file_manager.ReturnEntriesCheckedCount()))
}

func sortImages(logger *zap.Logger, cfg config, moveFile func(*zap.Logger, string, string) error) error {
	imageFiles, err := file_manager.GetFilesAllDepths(
		logger, cfg.sourcePath, image_manager.GetImageTypes(), true, image_manager.GetPhoto)
	if err != nil {
		return fmt.Errorf("failed to get image files from all depths: %w", err)
	}

	logger.Info("Got image files", zap.Int("count", len(imageFiles)))

	// sorting into folder structure of "<type>/<year>/<month>/<day>/<file>"
	// where type is either raw, edited or other,
	// other will be of format "<other>/<year>/<file>"
	usingImageFilesWithPath(logger, cfg, imageFiles, moveFile)
	return nil
}

func sortVideos(logger *zap.Logger, cfg config) error {
	err := video_manager.InitExifTool()
	if err != nil {
		return fmt.Errorf("failed to init exiftool: %w", err)
	}

	videoFiles, err := file_manager.GetFilesAllDepths(
		logger, cfg.sourcePath, video_manager.GetVideoTypes(), true, video_manager.GetVideo)
	if err != nil {
		return fmt.Errorf("failed to get video files from all depths: %w", err)
	}

	logger.Info("Got video files", zap.Int("count", len(videoFiles)))

	// sorting into folder structure of "<type>/<year>/<file>"
	// where type is either wildlife or other
	usingVideoFilesWithPath(logger, cfg, videoFiles)
	return nil
}

func usingSortedFolders(logger *zap.Logger, cfg config, imageFiles map[string]image_manager.ImageData) {
	sortedFolders := file_manager.SortFilesByDate(imageFiles, image_manager.GetTimestamp)

	logger.Debug("sorted files by date", zap.Any("sortedFolders", sortedFolders))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.destinationPath),
			zap.Error(err))
	}

	for folderName, files := range sortedFolders {
		err := file_manager.CreateFolderIfNotExists(logger, cfg.destinationPath+"/"+folderName)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", folderName),
				zap.Error(err))
		}

		for _, file := range files {
			err := file_manager.MoveAndRenameFile(
				logger,
				file.GetFilePath(),
				cfg.destinationPath+"/"+folderName+"/"+file.GetFileName())
			if err != nil {
				logger.Fatal("failed to copy and rename file",
					zap.String("destination", cfg.destinationPath+"/"+folderName+"/"+file.GetFileName()),
					zap.String("file", file.GetFileName()),
					zap.Error(err))
			}
		}
	}
}

func usingImageFilesWithPath(logger *zap.Logger, cfg config, imageFiles map[string]image_manager.ImageData,
	moveFile func(*zap.Logger, string, string) error) {
	logger.Info("Sorting files using source paths", zap.String("destinationPath", cfg.destinationPath))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.destinationPath),
			zap.Error(err))
	}

	filesWithPath := file_manager.AddFolderPathToFile(
		logger,
		imageFiles,
		sorting.AddingFolderToImagePath,
	)

	for _, file := range filesWithPath {
		logger.Debug("copying file",
			zap.String("destination", cfg.destinationPath+"/"+file.DestPath),
			zap.String("file", file.GetFileName()),
			zap.String("cameraModel", file.GetCameraModel()))

		err := file_manager.CreatePathFoldersIfDoesntExists(logger, cfg.destinationPath, file.DestPath)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", file.GetFilePath()),
				zap.Error(err))
		}

		err = moveFile(
			logger,
			file.GetFilePath(),
			cfg.destinationPath+"/"+file.DestPath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", testDestinationPath+"/"+file.DestPath),
				zap.String("file", file.GetFileName()),
				zap.Error(err))
		}
	}
}

func usingVideoFilesWithPath(logger *zap.Logger, cfg config, videoFiles map[string]video_manager.VideoData) {
	logger.Info("Sorting files using source paths", zap.String("destinationPath", cfg.destinationPath))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.destinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.destinationPath),
			zap.Error(err))
	}

	filesWithPath := file_manager.AddFolderPathToFile(
		logger,
		videoFiles,
		sorting.AddingFolderToVideoPath,
	)

	for _, file := range filesWithPath {
		logger.Debug("copying file",
			zap.String("destination", cfg.destinationPath+"/"+file.DestPath),
			zap.String("file", file.GetFileName()),
			zap.String("cameraModel", file.GetCameraModel()))

		err := file_manager.CreatePathFoldersIfDoesntExists(logger, cfg.destinationPath, file.DestPath)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", file.GetFilePath()),
				zap.Error(err))
		}

		err = file_manager.MoveAndRenameFile(
			logger,
			file.GetFilePath(),
			cfg.destinationPath+"/"+file.DestPath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", testDestinationPath+"/"+file.DestPath),
				zap.String("file", file.GetFileName()),
				zap.Error(err))
		}
	}
}

func nonRecognisedFileSorter(logger *zap.Logger) {
	nonRecognisedFiles, err := file_manager.GetFilesSingleFolder(
		logger,
		testSourcePath,
		image_manager.GetImageTypes(),
		false,
		func(path string) (string, error) { return path, nil },
	)

	err = file_manager.CreateFolderIfNotExists(logger, testNonImageFilesDestinationPath)
	if err != nil {
		logger.Fatal("failed to create folder in destination path",
			zap.String("folderName", testNonImageFilesDestinationPath),
			zap.Error(err))
	}

	for _, filePath := range nonRecognisedFiles {
		if strings.Contains(filePath, "DS_Store") {
			continue
		}
		err := file_manager.MoveAndRenameFile(
			logger,
			filePath,
			testDestinationPath+testNonImageFilesDestinationPath+filePath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", testDestinationPath+testNonImageFilesDestinationPath+filePath),
				zap.String("file", filePath),
				zap.Error(err))
		}
	}
}
