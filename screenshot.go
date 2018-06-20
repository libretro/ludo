package main

import (
	"image"
	"image/png"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/all-core/gl"
)

func screenshotName() string {
	name := filepath.Base(g.gamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	date := time.Now().Format("2006-01-02-15-04-05")
	return name + "@" + date + ".png"
}

func takeScreenshot() {
	usr, _ := user.Current()
	g.menuActive = false
	videoRender()
	fbw, fbh := window.GetFramebufferSize()
	img := image.NewNRGBA(image.Rect(0, 0, fbw, fbh))
	gl.ReadPixels(0, 0, int32(fbw), int32(fbh), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	fd, _ := os.Create(usr.HomeDir + "/.playthemall/screenshots/" + screenshotName())
	png.Encode(fd, imaging.FlipV(img))
	g.menuActive = true
}
