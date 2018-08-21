package input

import (
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/settings"
	"github.com/libretro/go-playthemall/state"
)

const (
	// ActionMenuToggle toggles the menu UI
	ActionMenuToggle uint32 = libretro.DeviceIDJoypadR3 + 1
	// ActionFullscreenToggle switches between fullscreen and windowed mode
	ActionFullscreenToggle uint32 = libretro.DeviceIDJoypadR3 + 2
	// ActionShouldClose will cause the program to shutdown
	ActionShouldClose uint32 = libretro.DeviceIDJoypadR3 + 3
	// ActionLast is used for iterating
	ActionLast uint32 = libretro.DeviceIDJoypadR3 + 4
)

// VideoInterface is used to pass a video instance to Process
type VideoInterface interface {
	Reconfigure(bool)
	SetShouldClose(bool)
}

// ProcessActions checks if certain keys are pressed and perform corresponding actions
func ProcessActions() {
	// Toggle the menu if ActionMenuToggle is pressed
	if Released[0][ActionMenuToggle] && state.Global.CoreRunning {
		state.Global.MenuActive = !state.Global.MenuActive
	}

	// Toggle fullscreen if ActionFullscreenToggle is pressed
	if Released[0][ActionFullscreenToggle] {
		settings.Settings.VideoFullscreen = !settings.Settings.VideoFullscreen
		vid.Reconfigure(settings.Settings.VideoFullscreen)
		// TODO context reset
		settings.Save()
	}

	// Close on escape
	if Pressed[0][ActionShouldClose] {
		vid.Window.SetShouldClose(true)
	}
}
