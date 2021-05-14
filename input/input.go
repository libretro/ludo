// Package input exposes the two input callbacks Poll and State needed by the
// libretro implementation. It uses GLFW to access keyboard and joypad, and
// takes care of binding and auto configuring joypads.
package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/video"
)

// MaxPlayers is the maximum number of players to poll input for
const MaxPlayers = 5

type joybinds map[bind]uint32

const btn = 0
const axis = 1

type bind struct {
	kind      uint32
	index     uint32
	direction float32
	threshold float32
}

// States can store the state of inputs for all players
type States [MaxPlayers][ActionLast]int16

// AnalogStates can store the state of analog inputs for all players
type AnalogStates [MaxPlayers][2]int16

// Input state for all the players
var (
	NewState States // input state for the current frame
	OldState States // input state for the previous frame
	Released States // keys just released during this frame
	Pressed  States // keys just pressed during this frame

	NewAnalogState AnalogStates // analog input state for the current frame
)

// Hot keys
const (
	// ActionMenuToggle toggles the menu UI
	ActionMenuToggle uint32 = libretro.DeviceIDJoypadR3 + 1
	// ActionFullscreenToggle switches between fullscreen and windowed mode
	ActionFullscreenToggle uint32 = libretro.DeviceIDJoypadR3 + 2
	// ActionShouldClose will cause the program to shutdown
	ActionShouldClose uint32 = libretro.DeviceIDJoypadR3 + 3
	// ActionFastForwardToggle will run the core as fast as possible
	ActionFastForwardToggle uint32 = libretro.DeviceIDJoypadR3 + 4
	// ActionLast is used for iterating
	ActionLast uint32 = libretro.DeviceIDJoypadR3 + 5
)

// joystickCallback is triggered when a joypad is plugged.
func joystickCallback(joy glfw.Joystick, event glfw.PeripheralEvent) {
	switch event {
	case glfw.Connected:
		if HasBinding(joy) {
			ntf.DisplayAndLog(ntf.Info, "Input", "Joystick #%d plugged: %s.", joy, glfw.Joystick.GetName(joy))
		} else {
			ntf.DisplayAndLog(ntf.Warning, "Input", "Joystick #%d plugged: %s but not configured.", joy, glfw.Joystick.GetName(joy))
		}
	case glfw.Disconnected:
		ntf.DisplayAndLog(ntf.Info, "Input", "Joystick #%d unplugged.", joy)
	default:
		ntf.DisplayAndLog(ntf.Warning, "Input", "Joystick #%d unhandled event: %d.", joy, event)
	}
}

var vid *video.Video

// Init initializes the input package
func Init(v *video.Video) {
	vid = v
	glfw.SetJoystickCallback(joystickCallback)
}

func floatToAnalog(v float32) int16 {
	return int16(v * 32767.0)
}

// pollJoypads process joypads of all players
func pollJoypads(state States, analogState AnalogStates) (States, AnalogStates) {
	for p := range state {
		buttonState := glfw.Joystick.GetButtons(glfw.Joystick(p))
		axisState := glfw.Joystick.GetAxes(glfw.Joystick(p))
		name := glfw.Joystick.GetName(glfw.Joystick(p))
		jb := joyBinds[name]
		if len(buttonState) > 0 {
			for k, v := range jb {
				switch k.kind {
				case btn:
					if int(k.index) < len(buttonState) &&
						glfw.Action(buttonState[k.index]) == glfw.Press {
						state[p][v] = 1
					}
				case axis:
					if int(k.index) < len(axisState) &&
						k.direction*axisState[k.index] > k.threshold*k.direction {
						state[p][v] = 1
					}
				}

				if settings.Current.MapAxisToDPad {
					if axisState[0] < -0.5 {
						state[p][libretro.DeviceIDJoypadLeft] = 1
					} else if axisState[0] > 0.5 {
						state[p][libretro.DeviceIDJoypadRight] = 1
					}
					if axisState[1] > 0.5 {
						state[p][libretro.DeviceIDJoypadDown] = 1
					} else if axisState[1] < -0.5 {
						state[p][libretro.DeviceIDJoypadUp] = 1
					}
				}
			}
		}
	}

	for p := range analogState {
		axisState := glfw.Joystick.GetAxes(glfw.Joystick(p))
		if len(axisState) >= 1 {
			analogState[p][0] = floatToAnalog(axisState[0])
			analogState[p][1] = floatToAnalog(axisState[1])
		}
	}

	return state, analogState
}

// pollKeyboard processes keyboard keys
func pollKeyboard(state States) States {
	for k, v := range keyBinds {
		if vid.Window.GetKey(k) == glfw.Press {
			state[0][v] = 1
		}
	}
	return state
}

// Compute the keys pressed or released during this frame
func getPressedReleased(new States, old States) (States, States) {
	for p := range new {
		for k := range new[p] {
			if new[p][k] == 1 && old[p][k] == 0 {
				Pressed[p][k] = 1
			} else {
				Pressed[p][k] = 0
			}
			if new[p][k] == 0 && old[p][k] == 1 {
				Released[p][k] = 1
			} else {
				Released[p][k] = 0
			}
		}
	}
	return Pressed, Released
}

// Poll calculates the input state. It is meant to be called for each frame.
func Poll() {
	NewState = States{}
	NewState, NewAnalogState = pollJoypads(NewState, NewAnalogState)
	NewState = pollKeyboard(NewState)
	Pressed, Released = getPressedReleased(NewState, OldState)

	// Store the old input state for comparisions
	OldState = NewState
}

// State is a callback passed to core.SetInputState
// It returns 1 if the button corresponding to the parameters is pressed
func State(port uint, device uint32, index uint, id uint) int16 {
	if port >= MaxPlayers {
		return 0
	}

	if device == libretro.DeviceJoypad {
		if id >= uint(ActionLast) || index > 0 {
			return 0
		}
		return NewState[port][id]
	}
	if device == libretro.DeviceAnalog {
		if id > uint(libretro.DeviceIDAnalogY) {
			// invalid
			return 0
		}

		switch uint32(index) {
		case libretro.DeviceIndexAnalogLeft:
			return NewAnalogState[port][id]
		case libretro.DeviceIndexAnalogRight:
			return NewAnalogState[port][id]
		}
	}

	return 0
}

// HasBinding returns true if the joystick has an autoconfig binding
func HasBinding(joy glfw.Joystick) bool {
	name := glfw.Joystick.GetName(joy)
	_, ok := joyBinds[name]
	return ok
}
