package menu

import (
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

var (
	repeatRight = withRepeat()
	repeatLeft  = withRepeat()
	repeatUp    = withRepeat()
	repeatDown  = withRepeat()
	repeatY     = withRepeat()
)

// Update takes care of calling the update method of the current scene.
// Each scene has it's own input logic to allow a variety of navigation systems.
func (m *Menu) Update(dt float32) {
	currentScene := m.stack[len(m.stack)-1]
	currentScene.update(dt)
}

// Used to increase scroll speed during long presses
func scrollSpeed(length float32) float32 {
	if length > 4 {
		return 0.005
	} else if length > 3 {
		return 0.01
	} else if length > 2 {
		return 0.02
	} else if length > 1 {
		return 0.04
	} else if length > 0.1 {
		return 0.08
	}
	return 0.15
}

// withRepeat wraps the logic that to allows firing repeated events when a key
// or button is hold. It is used mainly for scrolling, where the scroll speed
// increases with time. It's better to use 1 withRepeat per key, to achieve
// isolation.
func withRepeat() func(dt float32, pressed bool, f func()) {
	// A closure to store the values of these 3 vars across repeated calls
	var cooldown, length, delay float32
	return func(dt float32, pressed bool, f func()) {
		cooldown -= dt
		if pressed {
			if cooldown <= 0 {
				f()
				cooldown = delay
			}
			length += dt
		} else {
			length = 0
		}
		delay = scrollSpeed(length)
	}
}

// This is the generic menu input handler. It encapsulate the logic to scroll
// vertically in entry lists, and also respond to presses on OK and Cancel.
func genericInput(list *entry, dt float32) {
	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown] == 1, func() {
		list.ptr++
		if list.ptr >= len(list.children) {
			list.ptr = 0
		}
		audio.PlayEffect(audio.Effects["down"])
		genericAnimate(list)
	})

	// Up
	repeatUp(dt, input.NewState[0][libretro.DeviceIDJoypadUp] == 1, func() {
		list.ptr--
		if list.ptr < 0 {
			list.ptr = len(list.children) - 1
		}
		audio.PlayEffect(audio.Effects["up"])
		genericAnimate(list)
	})

	// Right
	if input.Released[0][libretro.DeviceIDJoypadRight] == 1 {
		if list.children[list.ptr].incr != nil {
			audio.PlayEffect(audio.Effects["up"])
			list.children[list.ptr].incr(1)
		}
	}

	// Left
	if input.Released[0][libretro.DeviceIDJoypadLeft] == 1 {
		if list.children[list.ptr].incr != nil {
			audio.PlayEffect(audio.Effects["down"])
			list.children[list.ptr].incr(-1)
		}
	}

	var confirmKey uint32
	var cancelKey uint32

	confirmKey = libretro.DeviceIDJoypadA
	cancelKey = libretro.DeviceIDJoypadB

	// Optionally swap confirm and cancel
	if settings.Current.SwapConfirm {
		confirmKey = libretro.DeviceIDJoypadB
		cancelKey = libretro.DeviceIDJoypadA
	}

	// OK
	if input.Released[0][confirmKey] == 1 {
		if list.children[list.ptr].callbackOK != nil {
			audio.PlayEffect(audio.Effects["ok"])
			list.children[list.ptr].callbackOK()
		}
	}

	// Cancel
	if input.Released[0][cancelKey] == 1 {
		if len(menu.stack) > 1 {
			audio.PlayEffect(audio.Effects["cancel"])
			menu.stack[len(menu.stack)-2].segueBack()
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}

	// X
	if input.Released[0][libretro.DeviceIDJoypadX] == 1 {
		if list.children[list.ptr].callbackX != nil {
			audio.PlayEffect(audio.Effects["ok"])
			list.children[list.ptr].callbackX()
		}
	}

	// Jump to next letter
	if input.Released[0][libretro.DeviceIDJoypadR] == 1 && len(list.indexes) > 0 {
		list.ptr = indexed(list, +1)
		audio.PlayEffect(audio.Effects["down"])
		genericAnimate(list)
	}

	// Jump to previous letter
	if input.Released[0][libretro.DeviceIDJoypadL] == 1 && len(list.indexes) > 0 {
		list.ptr = indexed(list, -1)
		audio.PlayEffect(audio.Effects["up"])
		genericAnimate(list)
	}
}

// indexed allows jumping directly to the next letter in playlists
func indexed(list *entry, offset int) int {
	curr := list.children[list.ptr].label[0]
	for i, t := range list.indexes {
		if curr == t.Char {
			if i+offset < 0 {
				return len(list.children) - 1
			}
			if i+offset > len(list.indexes)-1 {
				return 0
			}
			return list.indexes[i+offset].Index
		}
	}
	return 0
}

var combo1, combo2 int

// ProcessHotkeys checks if certain keys are pressed and perform corresponding actions
func (m *Menu) ProcessHotkeys() {
	// Disable all hot keys on the exit dialog
	currentScene := m.stack[len(m.stack)-1]
	if currentScene.Entry().label == "Confirm Dialog" {
		return
	}

	// First menu combo
	if input.NewState[0][libretro.DeviceIDJoypadL3] == 1 && input.NewState[0][libretro.DeviceIDJoypadR3] == 1 {
		combo1++
	} else {
		combo1 = 0
	}

	// Second menu combo
	if input.NewState[0][libretro.DeviceIDJoypadStart] == 1 && input.NewState[0][libretro.DeviceIDJoypadSelect] == 1 {
		combo2++
	} else {
		combo2 = 0
	}

	// Toggle the menu if ActionMenuToggle or the combo L3+R3 is pressed
	if (input.Pressed[0][input.ActionMenuToggle] == 1 || combo1 == 1 || combo2 == 1) && state.CoreRunning {
		state.MenuActive = !state.MenuActive
		state.FastForward = false
		if state.MenuActive {
			audio.PlayEffect(audio.Effects["notice"])
		} else {
			audio.PlayEffect(audio.Effects["notice_back"])
		}
	}

	// Toggle fullscreen if ActionFullscreenToggle is released
	if input.Released[0][input.ActionFullscreenToggle] == 1 {
		settings.Current.VideoFullscreen = !settings.Current.VideoFullscreen
		m.Reconfigure(settings.Current.VideoFullscreen)
		m.ContextReset()
		err := settings.Save()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", "Error saving settings: %s", err)
		}
	}

	if input.Pressed[0][input.ActionFastForwardToggle] == 1 && !state.MenuActive {
		state.FastForward = !state.FastForward
		if state.FastForward {
			ntf.DisplayAndLog(ntf.Info, "Menu", "Fast forward ON")
		} else {
			ntf.DisplayAndLog(ntf.Info, "Menu", "Fast forward OFF")
		}
	}

	// Close if ActionShouldClose is pressed, but display a confirmation dialog
	// in case a game is running
	if input.Pressed[0][input.ActionShouldClose] == 1 {
		askQuitConfirmation(func() {
			m.SetShouldClose(true)
		})
	}
}
