package image_manager

import (
	"fmt"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
)

type ImageData struct {
	filePath    string
	cameraModel string
	timestamp   time.Time
}

func toImageData(e exif2.Exif, path string) ImageData {
	return ImageData{
		filePath:    path,
		cameraModel: e.Model,
		timestamp:   e.DateTimeOriginal(),
	}
}

func GetTimestamp(i ImageData) time.Time {
	return i.timestamp
}

func GetPhoto(path string) (ImageData, error) {
	var i ImageData
	f, err := os.Open(path)
	if err != nil {
		return i, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		return i, fmt.Errorf("failed to decode image: %w", err)
	}

	return toImageData(e, path), nil
}
