package video

import (
	"image"
	"image/png"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/all-core/gl"
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
	usr, _ := user.Current()
	state.Global.MenuActive = false
	video.Render()
	fbw, fbh := video.Window.GetFramebufferSize()
	img := image.NewNRGBA(image.Rect(0, 0, fbw, fbh))
	gl.ReadPixels(0, 0, int32(fbw), int32(fbh), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	fd, _ := os.Create(usr.HomeDir + "/.ludo/screenshots/" + screenshotName())
	png.Encode(fd, imaging.FlipV(img))
	state.Global.MenuActive = true
}
