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
	VideoFullscreen   bool              `json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int               `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	AudioVolume       float32           `json:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles   bool              `json:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist   map[string]string `json:"core_for_playlist"`
	CoresDirectory    string            `json:"cores_dir" label:"Cores Directory" fmt:"%s" widget:"dir"`
}

// Current stores the current settings at runtime
var Current Settings

// Defaults stores default values for settings
var Defaults = Settings{
	VideoFullscreen:   false,
	VideoMonitorIndex: 0,
	AudioVolume:       0.5,
	ShowHiddenFiles:   true,
	CoreForPlaylist: map[string]string{
		"Atari - 2600":                                   "stella_libretro",
		"Atari - 5200":                                   "atari800_libretro",
		"Atari - 7800":                                   "prosystem_libretro",
		"Atari - Jaguar":                                 "virtualjaguar_libretro",
		"Atari - Lynx":                                   "handy_libretro",
		"Atari - ST":                                     "hatari_libretro",
		"Bandai - WonderSwan Color":                      "mednafen_wswan_libretro",
		"Bandai - WonderSwan":                            "mednafen_wswan_libretro",
		"Cave Story":                                     "nxengine_libretro",
		"ChaiLove":                                       "chailove_libretro",
		"Coleco - ColecoVision":                          "bluemsx_libretro",
		"FB Alpha - Arcade Games":                        "fbalpha_libretro",
		"GCE - Vectrex":                                  "vecx_libretro",
		"NEC - PC Engine SuperGrafx":                     "mednafen_supergrafx_libretro",
		"NEC - PC Engine - TurboGrafx 16":                "mednafen_pce_fast_libretro",
		"Nintendo - Game Boy Advance":                    "mgba_libretro",
		"Nintendo - Game Boy Color":                      "gambatte_libretro",
		"Nintendo - Game Boy":                            "gambatte_libretro",
		"Nintendo - Nintendo Entertainment System":       "nestopia_libretro",
		"Nintendo - Super Nintendo Entertainment System": "snes9x_libretro",
		"Nintendo - Virtual Boy":                         "mednafen_vb_libretro",
		"Sega - 32X":                                     "picodrive_libretro",
		"Sega - Game Gear":                               "genesis_plus_gx_libretro",
		"Sega - Master System - Mark III":                "genesis_plus_gx_libretro",
		"Sega - Mega Drive - Genesis":                    "genesis_plus_gx_libretro",
		"Sega - PICO":                                    "picodrive_libretro",
		"SNK - Neo Geo Pocket Color":                     "mednafen_ngp_libretro",
		"SNK - Neo Geo Pocket":                           "mednafen_ngp_libretro",
		"Sony - PlayStation":                             "mednafen_psx_libretro",
	},
	CoresDirectory: "./cores",
}

// Load loads settings from the home directory.
// If the settings file doesn't exists, it will return an error and
// set all the settings to their default value.
func Load() error {
	// Set default values for settings
	Current = Defaults

	usr, err := user.Current()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(usr.HomeDir + "/.ludo/settings.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &Current)

	Save()

	return err
}

// Save saves the current configuration to the home directory
func Save() error {
	usr, _ := user.Current()

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
	return "", errors.New("Default core not set")
}
