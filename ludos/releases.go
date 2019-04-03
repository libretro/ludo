package ludos

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"

	ntf "github.com/libretro/ludo/notifications"
)

// UpdatesDir is where releases should be saved to
const UpdatesDir = "/storage/.updates/"
const releasesURL = "https://api.github.com/repos/libretro/LudOS/releases"

var client = grab.NewClient()
var downloading bool
var libreELECArch = os.Getenv("LIBREELEC_ARCH")

// GHAsset is an asset attached to a github release
type GHAsset struct {
	Name               string
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GHRelease is a LudOS release hosted on github
type GHRelease struct {
	Name   string
	Assets []GHAsset
}

// GetReleases will get and decode the json from github api, and return the
// list of LudOS releases
func GetReleases() (*[]GHRelease, error) {
	r, err := http.Get(releasesURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	releases := []GHRelease{}
	err = json.NewDecoder(r.Body).Decode(&releases)
	return &releases, err
}

// FilterAssets finds and return the asset matching the LIBREELEC_ARCH
func FilterAssets(assets []GHAsset) *GHAsset {
	for _, asset := range assets {
		if strings.Contains(asset.Name, libreELECArch) {
			return &asset
		}
	}
	return nil
}

// DownloadRelease will download a LudOS release from github
func DownloadRelease(name, path, url string) {
	if downloading {
		ntf.DisplayAndLog(ntf.Error, "Menu", "A download is already in progress")
		return
	}

	nid := ntf.DisplayAndLog(ntf.Info, "Menu", "0/100 Downloading %s", name)
	downloading = true
	defer func() { downloading = false }()

	req, err := grab.NewRequest(path, url)
	if err != nil {
		ntf.Update(nid, ntf.Error, err.Error())
		return
	}

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			ntf.Update(nid, ntf.Info, "%.0f/100 Downloading %s", 100*resp.Progress(), name)

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		ntf.Update(nid, ntf.Error, err.Error())
		return
	}

	ntf.Update(nid, ntf.Success, "Done downloading. You can now reboot your system.")
}
