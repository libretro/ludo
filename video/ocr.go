package video

import (
	"fmt"
	"strings"
	"time"

	"github.com/libretro/ludo/ocr"
)

// OCRCurrentFrame captures the current frame and returns extracted text using the OCR engine.
func (video *Video) OCRCurrentFrame() (string, error) {
	img, err := video.CaptureFrameImage()
	if err != nil {
		return "", fmt.Errorf("capture frame: %w", err)
	}

	text, err := ocr.TextFromImage(img)
	if err != nil {
		return "", fmt.Errorf("run OCR: %w", err)
	}

	return text, nil
}

// SetSubtitle displays text at the bottom of the screen for a given duration.
func (video *Video) SetSubtitle(text string, duration time.Duration) {
	video.subtitleText = strings.TrimSpace(text)
	if video.subtitleText == "" {
		video.subtitleUntil = time.Time{}
		return
	}

	video.subtitleUntil = time.Now().Add(duration)
}
