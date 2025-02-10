package sorting

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/genutils"
)

var (
	rawFileTypes    = []string{"raw", "cr3", "cr2"}
	editedFileTypes = []string{"jpg", "jpeg"}

	editedFilesContainText = []string{"dxo", "enhanced", "cr3"}
	editedFilesSiffexes    = []string{"e", "tz", "ps", "nr"}
	editedDulicateSuffixes = []string{
		"_0%d", "_%d",
		"-0%d", "-%d",
		"(0%d)", "(%d)",
		"0%d", "%d",
		"~0%d", "~%d"}

	acceptedCameraModels = []string{"canon", "dc-fz82"}
)

// IsUsableFileType checks if the file type is in the list of file types, if includeFiles is true
// it will return true if the file type is in the list, if includeFiles is false it will return true
// if the file type is not in the list.
func IsUsableFileType(fileTypes []string, name string, includeFiles bool) bool {
	splitName := strings.Split(name, ".")
	if !(len(splitName) == 2) {
		return false
	}

	var result bool
	for _, fileType := range fileTypes {
		if strings.ToLower(fileType) == strings.ToLower(splitName[1]) && includeFiles {
			result = true
		}
	}
	if !includeFiles && !result {
		result = true
	}
	return result
}

func IsEditedOrRaw(logger *zap.Logger, i image_manager.ImageData) string {
	name := i.GetFileName()

	splitName := strings.Split(name, ".")
	if len(splitName) < 2 {
		return ""
	}

	fileType := strings.ToLower(splitName[1])
	// if the file is in the rawFileTypes it will only be a raw file
	if genutils.InArray(rawFileTypes, fileType) {
		logger.Debug("file is a raw file",
			zap.String("fileType", fileType))
		return "raw"
	}

	fileName := strings.ToLower(splitName[0])
	switch {
	case !genutils.InArray(editedFileTypes, fileType):
		logger.Debug("file is not an edited file type",
			zap.String("fileType", fileType))
		return "other"
	case isAcceptedCameraModel(i):
		logger.Debug("file has an accepted camera model",
			zap.String("cameraModel", i.GetCameraModel()))
		return "edited"
	case genutils.ContainsFromArray(editedFilesContainText, fileName):
		logger.Debug("file contains edited file text",
			zap.String("fileName", fileName))
		return "edited"
	case hasEditSuffix(name):
		logger.Debug("file has an edit suffix",
			zap.String("fileName", fileName))
		return "edited"
	}

	logger.Debug("file is not a raw or edited file",
		zap.String("fileType", fileType),
		zap.String("fileName", fileName))
	return "other"
}

func isAcceptedCameraModel(image image_manager.ImageData) bool {
	if genutils.StringsContainInArray(acceptedCameraModels, image.GetCameraModel()) {
		return true
	}
	return false
}

func hasEditSuffix(name string) bool {
	for i := 0; i < 10; i++ {
		for _, dupeSuffixFormat := range editedDulicateSuffixes {
			dupeSuffix := fmt.Sprintf(dupeSuffixFormat, i)
			if strings.HasSuffix(name, dupeSuffix) {
				for _, suffix := range editedFilesSiffexes {
					fullSuffix := dupeSuffix + suffix
					if strings.HasSuffix(name, fullSuffix) {
						return true
					}
				}
			}
		}
	}
	return false
}
