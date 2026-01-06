package video

import (
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

// renderScreenshot draws the core frame centered at core resolution and returns the pixel bounds to read back.
func (video *Video) renderScreenshot() (int32, int32, int32, int32) {
	fbw, fbh := video.Window.GetFramebufferSize()
	w := float32(video.width)
	h := float32(video.height)
	if w == 0 || h == 0 {
		w = float32(video.Geom.BaseWidth)
		h = float32(video.Geom.BaseHeight)
	}

	// Center the core-resolution quad in the framebuffer.
	x := (float32(fbw) - w) / 2
	y := (float32(fbh) - h) / 2

	va := video.vertexArray(x, y, w, h, 1.0)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	readX := int32(math.Round(float64(x)))
	readY := int32(math.Round(float64(fbh) - math.Round(float64(y+h))))
	readW := int32(math.Round(float64(w)))
	readH := int32(math.Round(float64(h)))

	bindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	return readX, readY, readW, readH
}

// CaptureFrameImage captures the output of video.Render and returns it as an image.
func (video *Video) CaptureFrameImage() (image.Image, error) {
	prevMenu := state.MenuActive
	state.MenuActive = false
	defer func() { state.MenuActive = prevMenu }()

	gl.UseProgram(video.defaultProgram)

	readX, readY, readW, readH := video.renderScreenshot()

	img := image.NewRGBA(image.Rect(0, 0, int(readW), int(readH)))

	gl.ReadPixels(
		readX, readY,
		readW, readH,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	gl.UseProgram(video.program)

	return imaging.FlipV(img), nil
}

// TakeScreenshot captures the output of video.Render and writes it to a file
func (video *Video) TakeScreenshot(name string) error {
	flipped, err := video.CaptureFrameImage()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(settings.Current.ScreenshotsDirectory, os.ModePerm); err != nil {
		return err
	}

	path := filepath.Join(settings.Current.ScreenshotsDirectory, name+".png")
	fd, err := os.Create(path)
	if err != nil {
		return err
	}

	return png.Encode(fd, flipped)
}
