package video

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/all-core/gl"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

func screenshotName() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	date := time.Now().Format("2006-01-02-15-04-05")
	return name + "@" + date + ".png"
}

func (video *Video) renderScreenshot() {
	avi := state.Global.Core.GetSystemAVInfo()
	video.Geom = avi.Geometry

	va := video.vertexArray(0, 0, float32(video.Geom.BaseWidth), float32(video.Geom.BaseHeight), 1.0)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	gl.UseProgram(video.program)
	video.updateMaskUniform()
	gl.Uniform4f(gl.GetUniformLocation(video.program, gl.Str("texColor\x00")), 1, 1, 1, 1)

	gl.BindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// TakeScreenshot captures the ouput of video.Render and writes it to a file
func (video *Video) TakeScreenshot() {
	_, fbh := video.Window.GetFramebufferSize()
	state.Global.MenuActive = false
	video.renderScreenshot()
	img := image.NewRGBA(image.Rect(0, 0, video.Geom.BaseWidth, video.Geom.BaseHeight))
	gl.ReadPixels(
		0, int32(fbh-video.Geom.BaseHeight),
		int32(video.Geom.BaseWidth), int32(video.Geom.BaseHeight),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	os.MkdirAll(settings.Current.ScreenshotsDirectory, os.ModePerm)
	path := filepath.Join(settings.Current.ScreenshotsDirectory, screenshotName())
	fd, _ := os.Create(path)
	flipped := imaging.FlipV(img)
	png.Encode(fd, flipped)
	state.Global.MenuActive = true
}
