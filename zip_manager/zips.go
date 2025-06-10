package zip_manager

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/photos-sorter/file_manager"
)

type ZipData struct {
	Name    string
	Path    string
	Entries int
}

func GetZipFiles(logger *zap.Logger, path string) (map[string]ZipData, error) {
	files, err := file_manager.GetFilesAllDepths[ZipData](logger, path, []string{".zip"}, true,
		func(logger *zap.Logger, filePath string) (ZipData, error) {
			return ZipData{
				Name: filepath.Base(filePath),
				Path: filePath,
			}, nil
		})
	if err != nil {
		return map[string]ZipData{}, fmt.Errorf("failed to get zip files: %w", err)
	}
	return files, nil
}

func UnzipFileFromZip(logger *zap.Logger, src, dst string) ([]string, error) {
	logger.Debug("getting file names from zip",
		zap.String("sourcePath", src))

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
