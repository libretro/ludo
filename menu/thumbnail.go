package menu

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/libretro/ludo/settings"

	"github.com/libretro/ludo/video"
)

func downloadThumbnail(list *entry, i int, url, folderPath, path string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}

	out, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		list.children[i].thumbnail = menu.icons["img-broken"]
		return
	}
	defer out.Close()

	io.Copy(out, resp.Body)
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

func drawThumbnail(list *entry, i int, system, gameName string, x, y, w, h, scale float32, color video.Color) {
	folderPath := filepath.Join(settings.Current.ThumbnailsDirectory, system, "Named_Snaps")
	path := filepath.Join(folderPath, gameName+".png")
	legalName := scrubIllegalChars(gameName)
	url := "http://thumbnails.libretro.com/" + system + "/Named_Snaps/" + legalName + ".png"

	if list.children[i].thumbnail == 0 || list.children[i].thumbnail == menu.icons["img-dl"] {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			list.children[i].thumbnail = video.NewImage(path)
		} else if list.children[i].thumbnail != menu.icons["img-dl"] {
			list.children[i].thumbnail = menu.icons["img-dl"]
			go downloadThumbnail(list, i, url, folderPath, path)
		}
	}

	vid.DrawImage(
		list.children[i].thumbnail,
		x, y, w, h, scale,
		color,
	)
}
