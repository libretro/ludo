package video

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/libretro/ludo/ocr"
)

// OCRCurrentFrame captures the current frame, preprocesses it for text, sends it to OpenAI, and displays the result as a subtitle.
func (video *Video) OCRCurrentFrame() error {
	img, err := video.CaptureFrameImage()
	if err != nil {
		return fmt.Errorf("capture frame: %w", err)
	}

	processed := ocr.Preprocess(img)
	path := filepath.Join(fmt.Sprintf("ludo-ocr-%d.png", time.Now().UnixNano()))
	if err := saveImage(path, processed); err != nil {
		log.Printf("[OCR] failed to save debug image: %v", err)
	} else {
		log.Printf("[OCR] saved debug image: %s", path)
	}

	go func(procImg image.Image) {
		text, err := ocr.TextFromImage(procImg)
		if err != nil {
			log.Printf("[OCR] failed to extract text: %v", err)
			return
		}
		video.SetSubtitle(text, 6*time.Second)
	}(processed)

	return nil
}

func saveImage(path string, img image.Image) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	return png.Encode(fd, img)
}
