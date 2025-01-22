package image_manager

import (
	"time"
)

type ImageData struct {
	FilePath    string
	CameraModel string
	Timestamp   time.Time
}
