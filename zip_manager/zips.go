package zip_manager

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"go.uber.org/zap"
)

func UnzipFileFromZip(logger *zap.Logger, src, dst string) ([]string, error) {
	logger.Debug("getting file names from zip",
		zap.String("sourcePath", src))

	archive, err := zip.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer archive.Close()

	for _, file := range archive.File {
		logger.Debug("found entry in zip",
			zap.String("fileName", file.Name))

		if file.FileInfo().IsDir() {
			logger.Debug("skipping directory",
				zap.String("fileName", file.Name))
			continue
		}
	}
	return nil, nil
}

func handleDirectory(file *zip.File) error {
	if !file.FileInfo().IsDir() {
		return fmt.Errorf("file is not a directory")
	}

	dirReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open directory: %w", err)
	}
	defer dirReader.Close()

	// Read the directory contents into a buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, dirReader)
	if err != nil {
		return fmt.Errorf("failed to read directory contents: %w", err)
	}

	dirArchive, err := zip.NewReader(bytes.NewReader(buf.Bytes()), file.FileInfo().Size())
	if err != nil {
		return fmt.Errorf("failed to read directory as zip: %w", err)
	}

	for _, f := range dirArchive.File {
		if f.FileInfo().IsDir() {
			err = handleDirectory(f)
			if err != nil {
				return fmt.Errorf("failed to handle directory: %w", err)
			}
		}
		fmt.Println("Found file in directory:", f.Name)
	}

	return nil
}
