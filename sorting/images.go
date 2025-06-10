package sorting

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/config"
)

func SortImages(logger *zap.Logger, cfg config.Config, moveFile func(*zap.Logger, string, string) error) error {
	imageFiles, err := file_manager.GetFilesAllDepths(
		logger, cfg.SourcePath, image_manager.GetImageTypes(), true, image_manager.GetPhoto)
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

func usingImageFilesWithPath(logger *zap.Logger, cfg config.Config,
	imageFiles map[string]image_manager.ImageData,
	moveFile func(*zap.Logger, string, string) error,
) {
	logger.Info("Sorting files using source paths", zap.String("destinationPath", cfg.DestinationPath))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.DestinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.DestinationPath),
			zap.Error(err))
	}

	filesWithPath := file_manager.AddFolderPathToFile(
		logger,
		imageFiles,
		addingFolderToImagePath,
	)

	file_manager.FilesToMoveCount = len(filesWithPath)
	for _, file := range filesWithPath {
		logger.Debug("copying/moving file",
			zap.String("destination", cfg.DestinationPath+"/"+file.DestPath),
			zap.String("file", file.GetFileName()),
			zap.String("cameraModel", file.GetCameraModel()))

		err := file_manager.CreatePathFoldersIfDoesntExists(logger, cfg.DestinationPath, file.DestPath)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", file.GetFilePath()),
				zap.Error(err))
		}

		err = moveFile(
			logger,
			file.GetFilePath(),
			cfg.DestinationPath+"/"+file.DestPath)
		if err != nil {
			logger.Fatal("failed to copy and rename file",
				zap.String("destination", cfg.DestinationPath+"/"+file.DestPath),
				zap.String("file", file.GetFileName()),
				zap.Error(err))
		}
	}
}
