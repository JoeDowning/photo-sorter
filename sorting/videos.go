package sorting

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/pkg/config"
	"github.com/photos-sorter/video_manager"
)

func SortVideos(logger *zap.Logger, cfg config.Config) error {
	err := video_manager.InitExifTool()
	if err != nil {
		return fmt.Errorf("failed to init exiftool: %w", err)
	}

	videoFiles, err := file_manager.GetFilesAllDepths(
		logger, cfg.SourcePath, video_manager.GetVideoTypes(), true, video_manager.GetVideo)
	if err != nil {
		return fmt.Errorf("failed to get video files from all depths: %w", err)
	}

	logger.Info("Got video files", zap.Int("count", len(videoFiles)))

	// sorting into folder structure of "<type>/<year>/<file>"
	// where type is either wildlife or other
	usingVideoFilesWithPath(logger, cfg, videoFiles)
	return nil
}

func usingVideoFilesWithPath(logger *zap.Logger, cfg config.Config, videoFiles map[string]video_manager.VideoData) {
	logger.Info("Sorting files using source paths", zap.String("destinationPath", cfg.DestinationPath))
	err := file_manager.CreateFolderIfNotExists(logger, cfg.DestinationPath)
	if err != nil {
		logger.Fatal("failed to create destination path",
			zap.String("destinationPath", cfg.DestinationPath),
			zap.Error(err))
	}

	filesWithPath := file_manager.AddFolderPathToFile(
		logger,
		videoFiles,
		addingFolderToVideoPath,
	)

	for _, file := range filesWithPath {
		logger.Debug("copying file",
			zap.String("destination", cfg.DestinationPath+"/"+file.DestPath),
			zap.String("file", file.GetFileName()),
			zap.String("cameraModel", file.GetCameraModel()))

		err := file_manager.CreatePathFoldersIfDoesntExists(logger, cfg.DestinationPath, file.DestPath)
		if err != nil {
			logger.Fatal("failed to create folder in destination path",
				zap.String("folderName", file.GetFilePath()),
				zap.Error(err))
		}

		err = file_manager.MoveAndRenameFile(
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
