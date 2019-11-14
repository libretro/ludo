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
	"os/user"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/libretro/ludo/ludos"
	"github.com/libretro/ludo/utils"
	"gopkg.in/yaml.v2"
)

// Settings is the list of available settings for the program. It serializes to YAML.
// Tags are used to set a human readable label and a format for the settings value.
// Widget sets the graphical representation of the value.
type Settings struct {
	VideoFullscreen   bool   `hide:"ludos" yaml:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int    `yaml:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	VideoFilter       string `yaml:"video_filter" label:"Video Filter" fmt:"<%s>"`

	AudioVolume     float32           `yaml:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	MenuAudioVolume float32           `yaml:"menu_audio_volume" label:"Menu Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles bool              `yaml:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
	CoreForPlaylist map[string]string `hide:"always" yaml:"core_for_playlist"`

	CoresDirectory       string `hide:"ludos" yaml:"cores_dir" label:"Cores Directory" fmt:"%s" widget:"dir"`
	AssetsDirectory      string `hide:"ludos" yaml:"assets_dir" label:"Assets Directory" fmt:"%s" widget:"dir"`
	DatabaseDirectory    string `hide:"ludos" yaml:"database_dir" label:"Database Directory" fmt:"%s" widget:"dir"`
	SavestatesDirectory  string `hide:"ludos" yaml:"savestates_dir" label:"Savestates Directory" fmt:"%s" widget:"dir"`
	SavefilesDirectory   string `hide:"ludos" yaml:"savefiles_dir" label:"Savefiles Directory" fmt:"%s" widget:"dir"`
	ScreenshotsDirectory string `hide:"ludos" yaml:"screenshots_dir" label:"Screenshots Directory" fmt:"%s" widget:"dir"`
	SystemDirectory      string `hide:"ludos" yaml:"system_dir" label:"System Directory" fmt:"%s" widget:"dir"`
	PlaylistsDirectory   string `hide:"ludos" yaml:"playlists_dir" label:"Playlists Directory" fmt:"%s" widget:"dir"`
	ThumbnailsDirectory  string `hide:"ludos" yaml:"thumbnail_dir" label:"Thumbnails Directory" fmt:"%s" widget:"dir"`

	SSHService       bool `hide:"app" yaml:"ssh_service" label:"SSH" widget:"switch" service:"sshd.service" path:"/storage/.cache/services/sshd.conf"`
	SambaService     bool `hide:"app" yaml:"samba_service" label:"Samba" widget:"switch" service:"smbd.service" path:"/storage/.cache/services/samba.conf"`
	BluetoothService bool `hide:"app" yaml:"bluetooth_service" label:"Bluetooth" widget:"switch" service:"bluetooth.service" path:"/storage/.cache/services/bluez.conf"`
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

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Set default values for settings
	Current = Defaults

	// If /etc/ludo.yml exists, override the defaults
	if _, err := os.Stat("/etc/ludo.yml"); !os.IsNotExist(err) {
		b, _ := ioutil.ReadFile("/etc/ludo.yml")
		err = yaml.Unmarshal(b, &Current)
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".ludo", "settings.yml"))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &Current)
	if err != nil {
		return err
	}

	// Those are special fields, their value is not saved in settings.yml but
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

	b, err := yaml.Marshal(Current)
	if err != nil {
		return err
	}

	fd, err := os.Create(filepath.Join(usr.HomeDir, ".ludo", "settings.yml"))
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
