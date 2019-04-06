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
var progress float64
var done bool
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
func DownloadRelease(path, url string) {
	if downloading {
		ntf.DisplayAndLog(ntf.Error, "Menu", "A download is already in progress")
		return
	}

	nid := ntf.DisplayAndLog(ntf.Info, "Menu", "Downloading update 0%%")
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
			ntf.Update(nid, ntf.Info, "Downloading update %.0f%%%% ", 100*resp.Progress())
			progress = resp.Progress()

		case <-resp.Done:
			// download is complete
			downloading = false
			done = true
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		ntf.Update(nid, ntf.Error, err.Error())
		downloading = false
		done = false
		return
	}

	ntf.Update(nid, ntf.Success, "Done downloading. You can now reboot your system.")
}

// IsDownloading returns true if the updater is currently downloading a release
func IsDownloading() bool {
	return downloading
}

// IsDone returns true when the download is finished
func IsDone() bool {
	return done
}

// GetProgress returns the download progress
func GetProgress() float64 {
	return progress
}
