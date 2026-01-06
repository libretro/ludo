package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

var (
	clientMu sync.Mutex
	client   *gosseract.Client
)

func initClient() *gosseract.Client {
	if client != nil {
		return client
	}

	client = gosseract.NewClient()
	client.SetLanguage("jpn")
	return client
}

// TextFromImage runs OCR on an image and returns the detected text.
func TextFromImage(img image.Image) (string, error) {
	processed := preprocess(img)

	path := filepath.Join(fmt.Sprintf("ludo-ocr-%d.png", time.Now().UnixNano()))
	if err := saveDebugImage(path, processed); err != nil {
		log.Printf("[OCR] failed to save debug image: %v", err)
	} else {
		log.Printf("[OCR] saved debug image: %s", path)
	}

	buf := bytes.Buffer{}
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	if err := enc.Encode(&buf, processed); err != nil {
		return "", err
	}

	clientMu.Lock()
	defer clientMu.Unlock()

	c := initClient()
	if err := c.SetImageFromBytes(buf.Bytes()); err != nil {
		return "", err
	}

	// Try a focused page segmentation first, then broaden if empty.
	c.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK)
	text, err := c.Text()
	if err == nil {
		text = strings.TrimSpace(text)
	}

	if err == nil && text == "" {
		c.SetPageSegMode(gosseract.PSM_AUTO)
		text, err = c.Text()
		if err == nil {
			text = strings.TrimSpace(text)
		}
	}

	return text, err
}

func preprocess(img image.Image) image.Image {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	// Binarize to keep only the observed text color (#cbc6dd).
	const (
		targetR = 0xcb
		targetG = 0xc6
		targetB = 0xdd
	)

	bin := image.NewGray(img.Bounds())
	for y := bin.Bounds().Min.Y; y < bin.Bounds().Max.Y; y++ {
		for x := bin.Bounds().Min.X; x < bin.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rr := int(r >> 8)
			gg := int(g >> 8)
			bb := int(b >> 8)

			if rr == targetR && gg == targetG && bb == targetB {
				bin.SetGray(x, y, color.Gray{Y: 255})
			} else {
				bin.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}


	// Upscale the frame to help OCR decode small text.
	scale := 6
	resized := imaging.Resize(bin, w*scale, h*scale, imaging.Linear)

	return resized
}

func saveDebugImage(path string, img image.Image) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	return png.Encode(fd, img)
}
