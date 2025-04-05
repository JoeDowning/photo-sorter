package video_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type VideoExifData struct {
	AudioBitsPerSample  int      `json:"AudioBitsPerSample"`
	AudioChannels       int      `json:"AudioChannels"`
	AudioFormat         string   `json:"AudioFormat"`
	AudioSampleRate     int      `json:"AudioSampleRate"`
	AvgBitrate          string   `json:"AvgBitrate"`
	Balance             int      `json:"Balance"`
	BitDepth            int      `json:"BitDepth"`
	Comment             string   `json:"Comment"`
	CompatibleBrands    []string `json:"CompatibleBrands"`
	CompressorID        string   `json:"CompressorID"`
	CreateDate          string   `json:"CreateDate"`
	CurrentTime         string   `json:"CurrentTime"`
	Directory           string   `json:"Directory"`
	Duration            string   `json:"Duration"`
	Encoder             string   `json:"Encoder"`
	ExifToolVersion     float64  `json:"ExifToolVersion"`
	FileAccessDate      string   `json:"FileAccessDate"`
	FileInodeChangeDate string   `json:"FileInodeChangeDate"`
	FileModifyDate      string   `json:"FileModifyDate"`
	FileName            string   `json:"FileName"`
	FilePermissions     string   `json:"FilePermissions"`
	FileSize            string   `json:"FileSize"`
	FileType            string   `json:"FileType"`
	FileTypeExtension   string   `json:"FileTypeExtension"`
	GraphicsMode        string   `json:"GraphicsMode"`
	HandlerDescription  string   `json:"HandlerDescription"`
	HandlerType         string   `json:"HandlerType"`
	HandlerVendorID     string   `json:"HandlerVendorID"`
	ImageHeight         int      `json:"ImageHeight"`
	ImageSize           string   `json:"ImageSize"`
	ImageWidth          int      `json:"ImageWidth"`
	MajorBrand          string   `json:"MajorBrand"`
	MatrixStructure     string   `json:"MatrixStructure"`
	MediaCreateDate     string   `json:"MediaCreateDate"`
	MediaDataOffset     int      `json:"MediaDataOffset"`
	MediaDataSize       float64  `json:"MediaDataSize"`
	MediaDuration       string   `json:"MediaDuration"`
	MediaHeaderVersion  int      `json:"MediaHeaderVersion"`
	MediaLanguageCode   string   `json:"MediaLanguageCode"`
	MediaModifyDate     string   `json:"MediaModifyDate"`
	MediaTimeScale      int      `json:"MediaTimeScale"`
	Megapixels          float64  `json:"Megapixels"`
	MIMEType            string   `json:"MIMEType"`
	MinorVersion        string   `json:"MinorVersion"`
	Model               string   `json:"Model"`
	ModifyDate          string   `json:"ModifyDate"`
	MovieHeaderVersion  int      `json:"MovieHeaderVersion"`
	NextTrackID         int      `json:"NextTrackID"`
	OpColor             string   `json:"OpColor"`
	PosterTime          string   `json:"PosterTime"`
	PreferredRate       int      `json:"PreferredRate"`
	PreferredVolume     string   `json:"PreferredVolume"`
	PreviewDuration     string   `json:"PreviewDuration"`
	PreviewTime         string   `json:"PreviewTime"`
	Rotation            int      `json:"Rotation"`
	SelectionDuration   string   `json:"SelectionDuration"`
	SelectionTime       string   `json:"SelectionTime"`
	SourceFile          string   `json:"SourceFile"`
	SourceImageHeight   int      `json:"SourceImageHeight"`
	SourceImageWidth    int      `json:"SourceImageWidth"`
	TimeScale           int      `json:"TimeScale"`
	Title               string   `json:"Title"`
	TrackCreateDate     string   `json:"TrackCreateDate"`
	TrackDuration       string   `json:"TrackDuration"`
	TrackHeaderVersion  int      `json:"TrackHeaderVersion"`
	TrackID             int      `json:"TrackID"`
	TrackLayer          int      `json:"TrackLayer"`
	TrackModifyDate     string   `json:"TrackModifyDate"`
	TrackVolume         string   `json:"TrackVolume"`
	VideoFrameRate      int      `json:"VideoFrameRate"`
	XResolution         int      `json:"XResolution"`
	YResolution         int      `json:"YResolution"`
}

func extractVideoDetails(logger *zap.Logger, raw map[string]interface{}) (VideoExifData, error) {
	for key, value := range raw {
		logger.Debug("raw data",
			zap.Any(key, value))
	}
	var data VideoExifData
	jsonData, err := json.Marshal(raw)
	if err != nil {
		return data, fmt.Errorf("failed to marshal raw data to JSON: %w", err)
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal JSON to VideoExifData: %w", err)
	}

	return data, nil
}

func exifDataToVideoData(logger *zap.Logger, data VideoExifData, path string) VideoData {
	camera := data.Model
	logger.Debug("camera model from exif data",
		zap.String("cameraModel", camera))
	if camera == "" {
		camera = data.Comment
		logger.Debug("camera model not found in exif data, using comment instead",
			zap.String("comment", data.Comment))
	}
	return VideoData{
		fileName:    data.FileName,
		filePath:    path,
		cameraModel: camera,
		timestamp:   parseTimestamp(data.CreateDate),
	}
}

func parseTimestamp(timestamp string) time.Time {
	t, err := time.Parse("2006:01:02 15:04:05", timestamp)
	if err != nil {
		return time.Time{}
	}
	return t
}
