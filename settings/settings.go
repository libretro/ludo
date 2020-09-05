// Package settings defines the app settings and functions to save and load
// those.
package settings

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/libretro/ludo/ludos"
	"github.com/libretro/ludo/utils"
	"github.com/pelletier/go-toml"
)

// Settings is the list of available settings for the program. It serializes to TOML.
// Tags are used to set a human readable label and a format for the settings value.
// Widget sets the graphical representation of the value.
type Settings struct {
	VideoFullscreen   bool   `hide:"ludos" toml:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int    `toml:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	VideoFilter       string `toml:"video_filter" label:"Video Filter" fmt:"<%s>"`
	VideoDarkMode     bool   `toml:"video_dark_mode" label:"Video Dark Mode" fmt:"%t" widget:"switch"`

	AudioVolume     float32           `toml:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	MenuAudioVolume float32           `toml:"menu_audio_volume" label:"Menu Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles bool              `toml:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist map[string]string `hide:"always" toml:"core_for_playlist"`

	CoresDirectory       string `hide:"ludos" toml:"cores_dir" label:"Cores Directory" fmt:"%s" widget:"dir"`
	AssetsDirectory      string `hide:"ludos" toml:"assets_dir" label:"Assets Directory" fmt:"%s" widget:"dir"`
	DatabaseDirectory    string `hide:"ludos" toml:"database_dir" label:"Database Directory" fmt:"%s" widget:"dir"`
	SavestatesDirectory  string `hide:"ludos" toml:"savestates_dir" label:"Savestates Directory" fmt:"%s" widget:"dir"`
	SavefilesDirectory   string `hide:"ludos" toml:"savefiles_dir" label:"Savefiles Directory" fmt:"%s" widget:"dir"`
	ScreenshotsDirectory string `hide:"ludos" toml:"screenshots_dir" label:"Screenshots Directory" fmt:"%s" widget:"dir"`
	SystemDirectory      string `hide:"ludos" toml:"system_dir" label:"System Directory" fmt:"%s" widget:"dir"`
	PlaylistsDirectory   string `hide:"ludos" toml:"playlists_dir" label:"Playlists Directory" fmt:"%s" widget:"dir"`
	ThumbnailsDirectory  string `hide:"ludos" toml:"thumbnail_dir" label:"Thumbnails Directory" fmt:"%s" widget:"dir"`

	SSHService       bool `hide:"app" toml:"ssh_service" label:"SSH" widget:"switch" service:"sshd.service" path:"/storage/.cache/services/sshd.conf"`
	SambaService     bool `hide:"app" toml:"samba_service" label:"Samba" widget:"switch" service:"smbd.service" path:"/storage/.cache/services/samba.conf"`
	BluetoothService bool `hide:"app" toml:"bluetooth_service" label:"Bluetooth" widget:"switch" service:"bluetooth.service" path:"/storage/.cache/services/bluez.conf"`
}

// Current stores the current settings at runtime
var Current Settings

// Defaults stores default values for settings
var Defaults = defaultSettings()

// Load loads settings from the home directory.
// If the settings file doesn't exists, it will return an error and
// set all the settings to their default value.
func Load() error {
	defer func() {
		err := Save()
		if err != nil {
			log.Println(err)
		}
	}()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Set default values for settings
	Current = Defaults

	// If /etc/ludo.toml exists, override the defaults
	if _, err := os.Stat("/etc/ludo.toml"); !os.IsNotExist(err) {
		b, _ := ioutil.ReadFile("/etc/ludo.toml")
		err = toml.Unmarshal(b, &Current)
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(filepath.Join(home, ".ludo", "settings.toml"))
	if err != nil {
		return err
	}
	err = toml.Unmarshal(b, &Current)
	if err != nil {
		return err
	}

	// Those are special fields, their value is not saved in settings.toml but
	// depends on the presence of some files
	ludos.InitializeServiceSettingsValues(structs.Fields(&Current))

	return nil
}

// Save saves the current configuration to the home directory
func Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(home, ".ludo"), os.ModePerm)
	if err != nil {
		return err
	}

	b, err := toml.Marshal(Current)
	if err != nil {
		return err
	}

	fd, err := os.Create(filepath.Join(home, ".ludo", "settings.toml"))
	if err != nil {
		return err
	}
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Println(err)
		}
	}()

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
