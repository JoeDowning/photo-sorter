package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	testSourcePath      = "/Users/joe.downing/Pictures/Photos/testing-folder/test-images"
	testZipSourcePath   = "/Users/joe.downing/Pictures/Photos/testing-folder/test-zips"
	testDestinationPath = "/Users/joe.downing/Pictures/Photos/testing-folder/sorted"

	seagateSourcePath      = "/Volumes/Seagate/takeouts"
	seagateDestinationPath = "/Volumes/Seagate new/sorted"

	rawDestinationPath    = "/Users/joe.downing/Pictures/Photos/toMove"
	editedDestinationPath = "/Users/joe.downing/Pictures/highResRedo"

	locationTest         = "test"
	locationTestZip      = "testZip"
	locationSeagate      = "seagate"
	locationBackupRaw    = "backupRaw"
	locationBackupEdited = "backupEdited"

	typeImages = "images"
	typeVideos = "videos"

	fileModeMove = "move"
	fileModeCopy = "copy"
)

var (
	testConfig = Config{
		SourcePath:      testSourcePath,
		DestinationPath: testDestinationPath,
	}
	testZipConfig = Config{
		SourcePath:      testZipSourcePath,
		DestinationPath: testDestinationPath,
	}
	seagateConfig = Config{
		SourcePath:      seagateSourcePath,
		DestinationPath: seagateDestinationPath,
	}
	backUpRawConfig = Config{
		SourcePath:      rawDestinationPath,
		DestinationPath: seagateDestinationPath,
	}
	backUpEditedConfig = Config{
		SourcePath:      editedDestinationPath,
		DestinationPath: seagateDestinationPath,
	}
)

type envConfig struct {
	FileType string `env:"file_type"`
	FileMode string `env:"file_mode"`
	Location string `env:"loc"`
	LogLevel string `env:"log"`
}

type Config struct {
	Mode            string
	IncludeZips     bool
	FileType        string
	FileMode        string
	SourcePath      string
	DestinationPath string
	LogLevel        string
}

func GetConfig() (Config, error) {
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get env config: %w", err)
	}

	var cfg Config
	cfg.LogLevel = envCfg.LogLevel
	switch envCfg.Location {
	case locationTest:
		cfg.DestinationPath = testConfig.DestinationPath
		cfg.SourcePath = testConfig.SourcePath
	case locationTestZip:
		cfg.DestinationPath = testZipConfig.DestinationPath
		cfg.SourcePath = testZipConfig.SourcePath
		cfg.IncludeZips = true
	case locationSeagate:
		cfg.DestinationPath = seagateConfig.DestinationPath
		cfg.SourcePath = seagateConfig.SourcePath
	case locationBackupRaw:
		cfg.DestinationPath = backUpRawConfig.DestinationPath
		cfg.SourcePath = backUpRawConfig.SourcePath
	case locationBackupEdited:
		cfg.DestinationPath = backUpEditedConfig.DestinationPath
		cfg.SourcePath = backUpEditedConfig.SourcePath
	default:
		return Config{}, fmt.Errorf("unknown location: %s (choices: %s, %s, %s)",
			envCfg.Location,
			locationTest,
			locationTestZip,
			locationSeagate)
	}

	switch envCfg.FileType {
	case typeImages:
		cfg.FileType = typeImages
	case typeVideos:
		cfg.FileType = typeVideos
	default:
		return Config{}, fmt.Errorf("unknown file type: %s (choices: %s, %s)",
			envCfg.FileType,
			typeImages,
			typeVideos)
	}

	switch envCfg.FileMode {
	case fileModeMove:
		cfg.FileMode = fileModeMove
	case fileModeCopy:
		cfg.FileMode = fileModeCopy
	default:
		return Config{}, fmt.Errorf("unknown file mode: %s (choices: %s, %s)",
			envCfg.FileMode,
			fileModeMove,
			fileModeCopy)
	}

	return cfg, nil
}
