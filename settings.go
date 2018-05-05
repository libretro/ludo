package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var settings struct {
	VideoScale      int     `json:"video_scale" label:"Video Scale"`
	VideoFullscreen bool    `json:"video_fullscreen" label:"Video Fullscreen"`
	AudioVolume     float64 `json:"audio_volume" label:"Audio Volume"`
}

func loadSettings() error {
	// Set default values
	settings.VideoScale = 3
	settings.VideoFullscreen = false
	settings.AudioVolume = 0.5

	b, err := slurp("settings.json")
	if err != nil {
		return err
	}
	json.Unmarshal(b, settings)
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
