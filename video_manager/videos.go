package video_manager

import (
	"fmt"
	"time"

	"github.com/barasher/go-exiftool"
	"go.uber.org/zap"
)

var (
	et              *exiftool.Exiftool
	NoFileInfoError = fmt.Errorf("no file info")
	videoFileTypes  = []string{"mp4", "mov", "avi"}

	cameraTypes = []string{"iphone", "pixel", "gardepro", "canon"}
)

type VideoData struct {
	fileName    string
	filePath    string
	cameraModel string
	timestamp   time.Time
	DestPath    string
}

func InitExifTool() error {
	var err error
	et, err = exiftool.NewExiftool()
	if err != nil {
		return fmt.Errorf("failed to initialize exiftool: %w", err)
	}
	return nil
}

func ClearupExifTool() {
	et.Close()
}

func (v VideoData) GetFileName() string {
	return v.fileName
}

func (v VideoData) GetFilePath() string {
	return v.filePath
}

func (v VideoData) GetCameraModel() string {
	return v.cameraModel
}

func GetTimestamp(v VideoData) time.Time {
	return v.timestamp
}

func GetVideo(logger *zap.Logger, path string) (VideoData, error) {
	var v VideoData
	fileInfos := et.ExtractMetadata(path)

	if fileInfos == nil {
		return v, NoFileInfoError
	}

	// check this for the cameraTypes to get that
	// confirm how to get the time data stamp if they are standard or not
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		exifData, err := extractVideoDetails(logger, fileInfo.Fields)
		if err != nil {
			logger.Error("failed to extract video details",
				zap.Error(err),
				zap.String("file", fileInfo.File),
				zap.Any("fields", fileInfo.Fields))
			return v, fmt.Errorf("failed to extract video details: %w", err)
		}
		v = exifDataToVideoData(logger, exifData, path)
	}

	return v, nil
}

func GetVideoTypes() []string {
	return videoFileTypes
}
