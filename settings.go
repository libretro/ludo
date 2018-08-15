package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/user"
	"sync"

	"github.com/fatih/structs"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/utils"
)

var settingsLock sync.Mutex

// settings is the list of available settings for the program.
// It serializes to JSON.
// Tags are used to set a human readable label and a format for the settings value.
var settings struct {
	VideoFullscreen   bool    `json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t" widget:"switch"`
	VideoMonitorIndex int     `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	AudioVolume       float32 `json:"audio_volume" label:"Audio Volume" fmt:"%.1f" widget:"range"`
	ShowHiddenFiles   bool    `json:"menu_showhiddenfiles" label:"Show Hidden Files" fmt:"%t" widget:"switch"`
}

type callbackIncrement func(*structs.Field, int)

// incrCallbacks is a map of callbacks called when a setting value
// is incremented or decremented
var incrCallbacks = map[string]callbackIncrement{
	"VideoFullscreen": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		videoConfigure(settings.VideoFullscreen)
		saveSettings()
	},
	"VideoMonitorIndex": func(f *structs.Field, direction int) {
		v := f.Value().(int)
		v += direction
		if v < 0 {
			v = 0
		}
		if v > len(glfw.GetMonitors())-1 {
			v = len(glfw.GetMonitors()) - 1
		}
		f.Set(v)
		videoConfigure(settings.VideoFullscreen)
		saveSettings()
	},
	"AudioVolume": func(f *structs.Field, direction int) {
		v := f.Value().(float32)
		v += 0.1 * float32(direction)
		f.Set(v)
		audioSetVolume(v)
		saveSettings()
	},
	"ShowHiddenFiles": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		saveSettings()
	},
}

// loadSettings loads settings from the home directory.
// If the settings file doesn't exists, it will return an error and
// set all the settings to their default value.
func loadSettings() error {
	settingsLock.Lock()
	defer settingsLock.Unlock()

	// Set default values
	settings.VideoFullscreen = false
	settings.VideoMonitorIndex = 0
	settings.AudioVolume = 0.5
	settings.ShowHiddenFiles = true

	usr, _ := user.Current()

	b, err := utils.Slurp(usr.HomeDir + "/.playthemall/settings.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &settings)
	return err
}

// saveSettings saves the current configuration to the home directory
func saveSettings() error {
	settingsLock.Lock()
	defer settingsLock.Unlock()

	usr, _ := user.Current()

	b, _ := json.MarshalIndent(settings, "", "  ")
	f, err := os.Create(usr.HomeDir + "/.playthemall/settings.json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}
