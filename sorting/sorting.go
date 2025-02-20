package sorting

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/photos-sorter/image_manager"
	"github.com/photos-sorter/pkg/genutils"
	"github.com/photos-sorter/video_manager"
)

var (
	rawImageTypes    = []string{"raw", "cr3", "cr2"}
	editedImageTypes = []string{"jpg", "jpeg"}

	editedImageContainText       = []string{"dxo", "enhanced", "cr3"}
	editedImageSuffixes          = []string{"e", "tz", "ps", "nr"}
	editedDuplicateImageSuffixes = []string{
		"_0%d", "_%d",
		"-0%d", "-%d",
		"(0%d)", "(%d)",
		"0%d", "%d",
		"~0%d", "~%d"}

	videoCameraModelKeywords = []string{"canon", "dc-fz82", "panasonic", "gardepro"}

	acceptedCameraModels = []string{"canon", "dc-fz82"}

	// fullPathFormat being: top folder / date sub folders / filename
	// the date sub folders are in the format: year / month / day /
	fullPathFormat = fmt.Sprintf("%s/%s%s")
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

func AddingFolderToImagePath(logger *zap.Logger, file image_manager.ImageData) image_manager.ImageData {
	editOrRawFile := isImageEditedOrRaw(logger, file)
	timestamp := image_manager.GetTimestamp(file)
	year := strconv.Itoa(timestamp.Year())
	if editOrRawFile == "other" {
		file.DestPath = createNewFullPath(editOrRawFile, createDateSubfolders(year), file.GetFileName())
	} else {
		month := strconv.Itoa(int(timestamp.Month()))
		if len(month) == 1 {
			month = "0" + month
		}
		day := strconv.Itoa(timestamp.Day())
		if len(day) == 1 {
			day = "0" + day
		}
		file.DestPath = createNewFullPath(editOrRawFile, createDateSubfolders(year, month, day), file.GetFileName())
	}

	return file
}

func AddingFolderToVideoPath(logger *zap.Logger, file video_manager.VideoData) video_manager.VideoData {
	timestamp := video_manager.GetTimestamp(file)
	year := strconv.Itoa(timestamp.Year())

	rootFolder := isVideoWildlifeOrNot(logger, file)

	file.DestPath = createNewFullPath(rootFolder, createDateSubfolders(year), file.GetFileName())
	return file
}

func isImageEditedOrRaw(logger *zap.Logger, i image_manager.ImageData) string {
	name := i.GetFileName()

	splitName := strings.Split(name, ".")
	if len(splitName) < 2 {
		return ""
	}

	fileName := strings.ToLower(splitName[0])
	fileType := strings.ToLower(splitName[1])
	logger = logger.With(
		zap.String("fileName", fileName),
		zap.String("fileType", fileType),
		zap.String("fullFileName", name))

	// if the file is in the rawFileTypes it will only be a raw file
	if genutils.InArray(rawImageTypes, fileType) {
		logger.Debug("file is a raw file")
		return "raw"
	}

	switch {
	case !genutils.InArray(editedImageTypes, fileType):
		logger.Debug("file is not an edited file type")
		return "other"
	case isAcceptedCameraModel(i):
		logger.Debug("file has an accepted camera model",
			zap.String("cameraModel", i.GetCameraModel()))
		return "edited"
	case genutils.ContainsFromArray(editedImageContainText, fileName):
		logger.Debug("file contains edited file text")
		return "edited"
	case hasEditSuffix(name):
		logger.Debug("file has an edit suffix")
		return "edited"
	default:
		logger.Debug("file is not a raw or edited file")
		return "other"
	}
}

func isVideoWildlifeOrNot(logger *zap.Logger, v video_manager.VideoData) string {
	name := v.GetFileName()
	camera := v.GetCameraModel()

	splitName := strings.Split(name, ".")
	if len(splitName) < 2 {
		return ""
	}

	fileName := strings.ToLower(splitName[0])
	logger = logger.With(
		zap.String("fileName", fileName),
		zap.String("fullFileName", name))

	//todo: add a check if the name has something useful in it?
	switch {
	case genutils.StringsContainInArray(videoCameraModelKeywords, camera):
		logger.Debug("video has a wildlife camera model",
			zap.String("cameraModel", camera))
		return "wildlife"
	default:
		logger.Debug("video does not have a wildlife camera model")
		return "other"
	}
}

func createDateSubfolders(subFolders ...string) string {
	var dateSubfolders string
	for _, subFolder := range subFolders {
		dateSubfolders += subFolder + "/"
	}
	return dateSubfolders
}

func createNewFullPath(topFolder, dateSubfolders, fileName string) string {
	return fmt.Sprintf(fullPathFormat, topFolder, dateSubfolders, fileName)
}

func isAcceptedCameraModel(image image_manager.ImageData) bool {
	if genutils.StringsContainInArray(acceptedCameraModels, image.GetCameraModel()) {
		return true
	}
	return false
}

func hasEditSuffix(name string) bool {
	for i := 0; i < 10; i++ {
		for _, dupeSuffixFormat := range editedDuplicateImageSuffixes {
			dupeSuffix := fmt.Sprintf(dupeSuffixFormat, i)
			if strings.HasSuffix(name, dupeSuffix) {
				for _, suffix := range editedImageSuffixes {
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
