package deskenv

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

const url = "https://api.github.com/repos/libretro/LudOS/releases"

// UpdatesDir is where releases should be saved to
const UpdatesDir = "/storage/.updates/"

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
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	releases := []GHRelease{}
	err = json.NewDecoder(r.Body).Decode(&releases)
	return &releases, err
}

// FilterAssets finds and return the asset matching the slug, slug can be
// Generic.x86_64 or RPi2.arm
func FilterAssets(slug string, assets []GHAsset) *GHAsset {
	for _, asset := range assets {
		if strings.Contains(asset.Name, slug) {
			return &asset
		}
	}
	return nil
}

// DownloadRelease will download a LudOS release from github
func DownloadRelease(filepath, url string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, r.Body)
	return err
}
