package video

import (
	"fmt"
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

// TakeScreenshot captures the ouput of video.Render and writes it to a file
func (video *Video) TakeScreenshot() {
	state.Global.MenuActive = false
	video.Render()
	fbWidth, fbHeight := video.Window.GetFramebufferSize()
	x, y, w, h := video.gameFrameQuad(fbWidth, fbHeight)
	fmt.Println(x, y, w, h)
	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	gl.ReadPixels(int32(x), int32(y), int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	os.MkdirAll(settings.Current.ScreenshotsDirectory, os.ModePerm)
	path := filepath.Join(settings.Current.ScreenshotsDirectory, screenshotName())
	fd, _ := os.Create(path)
	flipped := imaging.FlipV(img)
	resized := imaging.Resize(flipped, video.Geom.BaseWidth, video.Geom.BaseHeight, imaging.NearestNeighbor)
	png.Encode(fd, resized)
	state.Global.MenuActive = true
}
