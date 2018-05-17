package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var settings struct {
	VideoScale        int     `json:"video_scale" label:"Video Scale" fmt:"%dx"`
	VideoFullscreen   bool    `json:"video_fullscreen" label:"Video Fullscreen" fmt:"%t"`
	VideoMonitorIndex int     `json:"video_monitor_index" label:"Video Monitor Index" fmt:"%d"`
	AudioVolume       float32 `json:"audio_volume" label:"Audio Volume" fmt:"%.1f"`
}

type settingCallbackIncrement func(*structs.Field, int)

var incrCallbacks = map[string]settingCallbackIncrement{
	"VideoScale": func(f *structs.Field, direction int) {
		v := f.Value().(int)
		v += direction
		f.Set(v)
		videoConfigure(video.geom, settings.VideoFullscreen)
		saveSettings()
	},
	"VideoFullscreen": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		videoConfigure(video.geom, settings.VideoFullscreen)
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
		videoConfigure(video.geom, settings.VideoFullscreen)
		saveSettings()
	},
	"AudioVolume": func(f *structs.Field, direction int) {
		v := f.Value().(float32)
		v += 0.1 * float32(direction)
		f.Set(v)
		audio.source.SetGain(v)
		saveSettings()
	},
}

func loadSettings() error {
	// Set default values
	settings.VideoScale = 3
	settings.VideoFullscreen = false
	settings.VideoMonitorIndex = 0
	settings.AudioVolume = 0.5

	b, err := slurp("settings.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &settings)
	if err != nil {
		return err
	}
	return nil
}

func saveSettings() {
	b, _ := json.MarshalIndent(settings, "", "  ")
	f, err := os.OpenFile("settings.json", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write(b)
	f.Close()
}
