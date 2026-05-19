package video

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/libretro/ludo/settings"
)

func (video *Video) screenshotOutputSize() (int32, int32) {
	width := int32(video.Geom.BaseWidth)
	height := int32(video.Geom.BaseHeight)
	if width < 1 {
		width = video.width
	}
	if height < 1 {
		height = video.height
	}
	if width < 1 {
		width = video.presentW
	}
	if height < 1 {
		height = video.presentH
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return width, height
}

func (video *Video) captureBackbufferScreenshot(name string) error {
	outW, outH := video.screenshotOutputSize()

	x := video.presentX
	y := video.presentY
	w := video.presentW
	h := video.presentH
	if w < 1 || h < 1 {
		fbw, fbh := video.Window.GetFramebufferSize()
		x, y, w, h = int32(0), int32(0), int32(fbw), int32(fbh)
	}

	var oldReadBuffer int32
	gl.GetIntegerv(gl.READ_BUFFER, &oldReadBuffer)
	gl.ReadBuffer(gl.BACK)

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gl.Finish()
	gl.ReadPixels(x, y, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.ReadBuffer(uint32(oldReadBuffer))

	flipped := imaging.FlipV(img)
	if int32(flipped.Bounds().Dx()) != outW || int32(flipped.Bounds().Dy()) != outH {
		flipped = imaging.Resize(flipped, int(outW), int(outH), imaging.Linear)
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
	coreGLState := video.BeginFrontendFrame()
	defer video.EndFrontendFrame(coreGLState)

	video.Render()
	return video.captureBackbufferScreenshot(name)
}

// TakeScreenshotInFrontendFrame captures a screenshot while the caller already
// owns a frontend frame and has already rendered the current game frame.
func (video *Video) TakeScreenshotInFrontendFrame(name string) error {
	return video.captureBackbufferScreenshot(name)
}
