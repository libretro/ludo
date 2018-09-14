// Package settings defines the app settings and functions to save and load
// those.
package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/user"
	"sync"

	"github.com/libretro/go-playthemall/utils"
)

var lock sync.Mutex

// Settings is the list of available settings for the program. It serializes to JSON.
// Tags are used to set a human readable label and a format for the settings value.
var Settings struct {
	VideoFullscreen   bool              `json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int               `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	AudioVolume       float32           `json:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles   bool              `json:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist   map[string]string `json:"core_for_playlist"`
}

// Load loads settings from the home directory.
// If the settings file doesn't exists, it will return an error and
// set all the settings to their default value.
func Load() error {
	lock.Lock()
	defer lock.Unlock()

	// Set default values
	Settings.VideoFullscreen = false
	Settings.VideoMonitorIndex = 0
	Settings.AudioVolume = 0.5
	Settings.ShowHiddenFiles = true
	Settings.CoreForPlaylist = map[string]string{
		"Nintendo - Super Nintendo Entertainment System": "snes9x_libretro",
		"Sega - Master System - Mark III":                "genesis_plus_gx_libretro",
		"Sega - Mega Drive - Genesis":                    "genesis_plus_gx_libretro",
	}

	usr, _ := user.Current()

	b, err := utils.Slurp(usr.HomeDir + "/.playthemall/settings.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &Settings)
	return err
}

// Save saves the current configuration to the home directory
func Save() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	b, _ := json.MarshalIndent(Settings, "", "  ")
	f, err := os.Create(usr.HomeDir + "/.playthemall/settings.json")
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
	usr, _ := user.Current()
	coresPath := usr.HomeDir + "/.playthemall/cores/"
	c := Settings.CoreForPlaylist[playlist]
	if c != "" {
		return coresPath + c + utils.CoreExt(), nil
	}
	return "", errors.New("Default core not set")
}
