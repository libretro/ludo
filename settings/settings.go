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

	"github.com/fatih/structs"
	"github.com/libretro/ludo/ludos"
	"github.com/libretro/ludo/utils"
)

// Settings is the list of available settings for the program. It serializes to JSON.
// Tags are used to set a human readable label and a format for the settings value.
// Widget sets the graphical representation of the value.
type Settings struct {
	VideoFullscreen   bool   `hide:"ludos" json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int    `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	VideoFilter       string `json:"video_filter" label:"Video Filter" fmt:"<%s>"`

	GLVersion       string            `hide:"always" json:"video_gl_version"`
	AudioVolume     float32           `json:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles bool              `json:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist map[string]string `hide:"always" json:"core_for_playlist"`

	CoresDirectory       string `hide:"ludos" json:"cores_dir" label:"Cores Directory" fmt:"%s" widget:"dir"`
	AssetsDirectory      string `hide:"ludos" json:"assets_dir" label:"Assets Directory" fmt:"%s" widget:"dir"`
	DatabaseDirectory    string `hide:"ludos" json:"database_dir" label:"Database Directory" fmt:"%s" widget:"dir"`
	SavestatesDirectory  string `hide:"ludos" json:"savestates_dir" label:"Savestates Directory" fmt:"%s" widget:"dir"`
	SavefilesDirectory   string `hide:"ludos" json:"savefiles_dir" label:"Savefiles Directory" fmt:"%s" widget:"dir"`
	ScreenshotsDirectory string `hide:"ludos" json:"screenshots_dir" label:"Screenshots Directory" fmt:"%s" widget:"dir"`
	SystemDirectory      string `hide:"ludos" json:"system_dir" label:"System Directory" fmt:"%s" widget:"dir"`
	PlaylistsDirectory   string `hide:"ludos" json:"playlists_dir" label:"Playlists Directory" fmt:"%s" widget:"dir"`
	ThumbnailsDirectory  string `hide:"ludos" json:"thumbnail_dir" label:"Thumbnails Directory" fmt:"%s" widget:"dir"`

	SSHService       bool `hide:"app" json:"ssh_service" label:"SSH" widget:"switch" service:"sshd.service" path:"/storage/.cache/services/sshd.conf"`
	SambaService     bool `hide:"app" json:"samba_service" label:"Samba" widget:"switch" service:"smbd.service" path:"/storage/.cache/services/samba.conf"`
	BluetoothService bool `hide:"app" json:"bluetooth_service" label:"Bluetooth" widget:"switch" service:"bluetooth.service" path:"/storage/.cache/services/bluez.conf"`
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

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Set default values for settings
	Current = Defaults

	// If /etc/ludo.json exists, override the defaults
	if _, err := os.Stat("/etc/ludo.json"); !os.IsNotExist(err) {
		b, _ := ioutil.ReadFile("/etc/ludo.json")
		json.Unmarshal(b, &Current)
	}

	b, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".ludo", "settings.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &Current)
	if err != nil {
		return err
	}

	// Those are special fields, their value is not saved in settings.json but
	// depends on the presence of some files
	ludos.InitializeServiceSettingsValues(structs.Fields(&Current))

	return nil
}

// Save saves the current configuration to the home directory
func Save() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".ludo"), os.ModePerm)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(Current, "", "  ")
	if err != nil {
		return err
	}

	fd, err := os.Create(filepath.Join(usr.HomeDir, ".ludo", "settings.json"))
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return fd.Sync()
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
