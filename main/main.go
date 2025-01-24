package main

import (
	"fmt"

	"github.com/photos-sorter/file_manager"
	"github.com/photos-sorter/image_manager"
)

var (
	fileTypes  = []string{"jpg", "jpeg", "raw", "cr3", "cr2", "png"}
	sourcePath = "/Users/joe.downing/Pictures/Photos/testing-folder"
	//destinationPath = ""
)

func main() {
	files, err := file_manager.GetFiles(sourcePath, fileTypes, image_manager.GetPhoto)
	if err != nil {
		panic(fmt.Errorf("failed to get files: %w", err))
	}

	sortedFolders := file_manager.SortFilesByDate(files, image_manager.GetTimestamp)

	// move files to destinationPath
	// change file name to have timestamp at the beginning of filename
	// creating folder per date if not existing

	fmt.Printf("Files: %v\n", files)
	fmt.Printf("Sorted Folders: %v\n", sortedFolders)
}
