package main

import (
	"fmt"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
)

var (
	imageFileTypes = []string{"jpg", "jpeg", "raw", "cr3", "cr2", "png"}
	//videoFileTypes  = []string{"mp4", "mov", "avi", "mkv", "gif"}
	sourcePath      = "/Users/joe.downing/Pictures/Photos/testing-folder"
	destinationPath = "/Users/joe.downing/Pictures/Photos/testing-folder/sorted/"
)

func main() {
	imageFiles, err := file_manager.GetFiles(sourcePath, imageFileTypes, image_manager.GetPhoto)
	if err != nil {
		panic(fmt.Errorf("failed to get files: %w", err))
	}

	//todo: add video files

	//todo: change below to have multiple file types and getTimestamps

	sortedFolders := file_manager.SortFilesByDate(imageFiles, image_manager.GetTimestamp)

	for folderName, files := range sortedFolders {

		//todo: change so that it is folders of year, month, date

		err := file_manager.CreateFolderIfNotExists(destinationPath + folderName)
		if err != nil {
			panic(fmt.Errorf("failed to create folder: %w", err))
		}

		for _, file := range files {
			err := file_manager.CopyAndRenameFile(
				file.GetFilePath(),
				destinationPath+folderName+"/"+file.GetFileName())
			if err != nil {
				panic(fmt.Errorf("failed to copy and rename file: %w", err))
			}
		}
	}
}
