// Package settings defines the app settings and functions to save and load
// those.
package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/libretro/ludo/utils"
)

// Settings is the list of available settings for the program. It serializes to JSON.
// Tags are used to set a human readable label and a format for the settings value.
// Widget sets the graphical representation of the value.
type Settings struct {
	VideoFullscreen      bool              `json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex    int               `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	GLVersion            string            `json:"video_gl_version"`
	AudioVolume          float32           `json:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles      bool              `json:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist      map[string]string `json:"core_for_playlist"`
	CoresDirectory       string            `json:"cores_dir" label:"Cores Directory" fmt:"%s" widget:"dir"`
	AssetsDirectory      string            `json:"assets_dir" label:"Assets Directory" fmt:"%s" widget:"dir"`
	DatabaseDirectory    string            `json:"database_dir" label:"Database Directory" fmt:"%s" widget:"dir"`
	SavestatesDirectory  string            `json:"savestates_dir" label:"Savestates Directory" fmt:"%s" widget:"dir"`
	SavefilesDirectory   string            `json:"savefiles_dir" label:"Savefiles Directory" fmt:"%s" widget:"dir"`
	ScreenshotsDirectory string            `json:"screenshots_dir" label:"Screenshots Directory" fmt:"%s" widget:"dir"`
	SystemDirectory      string            `json:"system_dir" label:"System Directory" fmt:"%s" widget:"dir"`
	PlaylistsDirectory   string            `json:"playlists_dir" label:"Playlists Directory" fmt:"%s" widget:"dir"`
	ThumbnailsDirectory  string            `json:"thumbnail_dir" label:"Thumbnails Directory" fmt:"%s" widget:"dir"`
}

// Current stores the current settings at runtime
var Current Settings

// Defaults stores default values for settings
var Defaults = defaultSettings()

// Load loads settings from the home directory.
// If the settings file doesn't exists, it will return an error and
// set all the settings to their default value.
func Load() error {
	defer Save()

	usr, _ := user.Current()

	// Set default values for settings
	Current = Defaults

	// If /etc/ludo.json exists, override the defaults
	if _, err := os.Stat("/etc/ludo.json"); !os.IsNotExist(err) {
		b, _ := ioutil.ReadFile("/etc/ludo.json")
		json.Unmarshal(b, &Current)
	}

	b, err := ioutil.ReadFile(usr.HomeDir + "/.ludo/settings.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &Current)

	return err
}

// Save saves the current configuration to the home directory
func Save() error {
	usr, _ := user.Current()

	os.MkdirAll(usr.HomeDir+"/.ludo", os.ModePerm)

	b, _ := json.MarshalIndent(Current, "", "  ")
	f, err := os.Create(usr.HomeDir + "/.ludo/settings.json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}

// CoreForPlaylist returns the absolute path of the default libretro core for
// a given playlist
func CoreForPlaylist(playlist string) (string, error) {
	c := Current.CoreForPlaylist[playlist]
	if c != "" {
		return filepath.Join(Current.CoresDirectory, c+utils.CoreExt()), nil
	}
	return "", errors.New("default core not set")
}
