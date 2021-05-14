package menu

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/video"
)

// Downloads a thumbnail from the web and cache it to the local filesystem.
func downloadThumbnail(list *entry, i int, url, folderPath, path string) {
	resp, err := http.Get(url)
	if err != nil {
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}

	imgFile, err := os.Create(path)
	if err != nil {
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}
	defer imgFile.Close()

	_, err = io.Copy(imgFile, resp.Body)
	if err != nil {
		list.children[i].thumbnail = menu.icons["img-broken"]
	}
}

// Scrub characters that are not cross-platform and/or violate the
// No-Intro filename standard.
func scrubIllegalChars(str string) string {
	str = strings.Replace(str, "&", "_", -1)
	str = strings.Replace(str, "*", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "`", "_", -1)
	str = strings.Replace(str, "<", "_", -1)
	str = strings.Replace(str, ">", "_", -1)
	str = strings.Replace(str, "?", "_", -1)
	str = strings.Replace(str, "|", "_", -1)
	return str
}

// Draws a thumbnail in the playlist scene.
func drawThumbnail(list *entry, i int, system, gameName string, x, y, w, h, scale float32, color video.Color) {
	folderPath := filepath.Join(settings.Current.ThumbnailsDirectory, system, "Named_Snaps")
	legalName := scrubIllegalChars(gameName)
	path := filepath.Join(folderPath, legalName+".png")
	url := "http://thumbnails.libretro.com/" + system + "/Named_Snaps/" + legalName + ".png"

	if list.children[i].thumbnail == 0 || list.children[i].thumbnail == menu.icons["img-dl"] {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			list.children[i].thumbnail = video.NewImage(path)
		} else if list.children[i].thumbnail != menu.icons["img-dl"] {
			list.children[i].thumbnail = menu.icons["img-dl"]
			go downloadThumbnail(list, i, url, folderPath, path)
		}
	}

	menu.DrawImage(
		list.children[i].thumbnail,
		x, y, w, h, scale,
		color,
	)
}

// Draws a thumbnail in the savestates scene.
func drawSavestateThumbnail(list *entry, i int, path string, x, y, w, h, scale float32, color video.Color) {
	if list.children[i].thumbnail == 0 {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			list.children[i].thumbnail = video.NewImage(path)
		}
	}

	menu.DrawImage(
		list.children[i].thumbnail,
		x, y, w, h, scale,
		color,
	)
}

func freeThumbnail(list *entry, i int) {
	if list.children[i].thumbnail != 0 &&
		list.children[i].thumbnail != menu.icons["img-dl"] &&
		list.children[i].thumbnail != menu.icons["img-broken"] {
		gl.DeleteTextures(1, &list.children[i].thumbnail)
		list.children[i].thumbnail = 0
	}
}
