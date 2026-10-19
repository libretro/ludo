package video

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

func (video *Video) captureBackbufferScreenshot(name string) error {
	outW, outH := video.Geom.BaseWidth, video.Geom.BaseHeight
	if outW < 1 || outH < 1 {
		return errors.New("invalid screenshot output size")
	}

	x := video.presentX
	y := video.presentY
	w := video.presentW
	h := video.presentH
	if !state.CoreRunning || video.width < 1 || video.height < 1 || w < 1 || h < 1 {
		return errors.New("no game frame to capture")
	}

	var oldReadBuffer int32
	gl.GetIntegerv(gl.READ_BUFFER, &oldReadBuffer)
	gl.ReadBuffer(gl.BACK)

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gl.Finish()
	gl.ReadPixels(x, y, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.ReadBuffer(uint32(oldReadBuffer))

	flipped := imaging.FlipV(img)
	if flipped.Bounds().Dx() != outW || flipped.Bounds().Dy() != outH {
		flipped = imaging.Resize(flipped, outW, outH, imaging.Linear)
	}

	err := os.MkdirAll(settings.Current.ScreenshotsDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(settings.Current.ScreenshotsDirectory, name+".png")
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, flipped)
}

// TakeScreenshot captures a screenshot without advancing the core.
func (video *Video) TakeScreenshot(name string) error {
	video.BeginFrontendFrame()

	video.Render()
	return video.captureBackbufferScreenshot(name)
}

// TakeScreenshotInFrontendFrame captures a screenshot while the caller already
// owns a frontend frame and has already rendered the current game frame.
func (video *Video) TakeScreenshotInFrontendFrame(name string) error {
	return video.captureBackbufferScreenshot(name)
}
