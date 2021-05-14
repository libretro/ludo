package video

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

// During the TakeScreenshot step, we need to render the current game frame at
// the right resolution to later capture it using ReadPixels. renderScreenshot
// taking care of this.
func (video *Video) renderScreenshot() {
	va := video.vertexArray(0, 0, float32(video.Geom.BaseWidth), float32(video.Geom.BaseHeight), 1.0)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	bindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// TakeScreenshot captures the ouput of video.Render and writes it to a file
func (video *Video) TakeScreenshot(name string) error {
	state.MenuActive = false
	defer func() { state.MenuActive = true }()

	gl.UseProgram(video.defaultProgram)

	video.renderScreenshot()

	img := image.NewRGBA(image.Rect(0, 0, video.Geom.BaseWidth, video.Geom.BaseHeight))

	_, fbh := video.Window.GetFramebufferSize()

	gl.ReadPixels(
		0, int32(fbh-video.Geom.BaseHeight),
		int32(video.Geom.BaseWidth), int32(video.Geom.BaseHeight),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	gl.UseProgram(video.program)

	err := os.MkdirAll(settings.Current.ScreenshotsDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(settings.Current.ScreenshotsDirectory, name+".png")
	fd, err := os.Create(path)
	if err != nil {
		return err
	}

	flipped := imaging.FlipV(img)

	return png.Encode(fd, flipped)
}
