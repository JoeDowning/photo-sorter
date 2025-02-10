package image_manager

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"

	"github.com/photos-sorter/pkg/genutils"
)

type ImageData struct {
	fileName    string
	filePath    string
	cameraModel string
	timestamp   time.Time
	DestPath    string
}

func toImageData(e exif2.Exif, name, path string) ImageData {
	h, m, s := e.DateTimeOriginal().Clock()
	prefix := fmt.Sprintf("%s%s%s_",
		genutils.PrefixZeros(2, strconv.Itoa(h)),
		genutils.PrefixZeros(2, strconv.Itoa(m)),
		genutils.PrefixZeros(2, strconv.Itoa(s)))
	return ImageData{
		fileName:    prefix + name,
		filePath:    path,
		cameraModel: e.Model,
		timestamp:   e.DateTimeOriginal(),
	}
}

func (i ImageData) GetFileName() string {
	return i.fileName
}

func (i ImageData) GetFilePath() string {
	return i.filePath
}

func (i ImageData) GetCameraModel() string {
	return strings.ToLower(i.cameraModel)
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

	sepPath := strings.Split(path, "/")
	return toImageData(e, sepPath[len(sepPath)-1], path), nil
}
