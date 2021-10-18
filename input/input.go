// Package input exposes the two input callbacks Poll and State needed by the
// libretro implementation. It uses GLFW to access keyboard and joypad, and
// takes care of binding and auto configuring joypads.
package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	lr "github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/video"
)

// MaxPlayers is the maximum number of players to poll input for
const MaxPlayers = 5

// States can store the state of inputs for all players
type States [MaxPlayers][ActionLast]int16

// AnalogStates can store the state of analog inputs for all players
type AnalogStates [MaxPlayers][2][2]int16

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
	ActionMenuToggle uint32 = lr.DeviceIDJoypadR3 + 1
	// ActionFullscreenToggle switches between fullscreen and windowed mode
	ActionFullscreenToggle uint32 = lr.DeviceIDJoypadR3 + 2
	// ActionShouldClose will cause the program to shutdown
	ActionShouldClose uint32 = lr.DeviceIDJoypadR3 + 3
	// ActionFastForwardToggle will run the core as fast as possible
	ActionFastForwardToggle uint32 = lr.DeviceIDJoypadR3 + 4
	// ActionLast is used for iterating
	ActionLast uint32 = lr.DeviceIDJoypadR3 + 5
)

// joystickCallback is triggered when a joypad is plugged.
func joystickCallback(joy glfw.Joystick, event glfw.PeripheralEvent) {
	switch event {
	case glfw.Connected:
		if joy.IsGamepad() {
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
	if !glfw.UpdateGamepadMappings(mappings) {
		log.Println("Failed to update mappings")
	}
	glfw.SetJoystickCallback(joystickCallback)
}

func floatToAnalog(v float32) int16 {
	return int16(v * 32767.0)
}

// pollJoypads process joypads of all players
func pollJoypads(state States, analogState AnalogStates) (States, AnalogStates) {
	p := 0
	for joy := glfw.Joystick(0); joy < glfw.JoystickLast; joy++ {
		if !joy.IsGamepad() {
			continue
		}
		pad := joy.GetGamepadState()
		if pad == nil {
			continue
		}

		// mapping pad buttons
		for k, v := range joyBinds {
			if pad.Buttons[k] == glfw.Press {
				state[p][v] = 1
			}
		}

		// mapping pad triggers
		if pad.Axes[glfw.AxisLeftTrigger] > 0.5 {
			state[p][lr.DeviceIDJoypadL2] = 1
		}
		if pad.Axes[glfw.AxisRightTrigger] > 0.5 {
			state[p][lr.DeviceIDJoypadR2] = 1
		}

		// mapping analog sticks
		analogState[p][lr.DeviceIndexAnalogLeft][lr.DeviceIDAnalogX] = floatToAnalog(pad.Axes[glfw.AxisLeftX])
		analogState[p][lr.DeviceIndexAnalogLeft][lr.DeviceIDAnalogY] = floatToAnalog(pad.Axes[glfw.AxisLeftY])
		analogState[p][lr.DeviceIndexAnalogRight][lr.DeviceIDAnalogX] = floatToAnalog(pad.Axes[glfw.AxisRightX])
		analogState[p][lr.DeviceIndexAnalogRight][lr.DeviceIDAnalogY] = floatToAnalog(pad.Axes[glfw.AxisRightY])

		// optionally mapping analog sticks to dpad
		if settings.Current.MapAxisToDPad {
			if pad.Axes[glfw.AxisLeftX] < -0.5 {
				state[p][lr.DeviceIDJoypadLeft] = 1
			} else if pad.Axes[glfw.AxisLeftX] > 0.5 {
				state[p][lr.DeviceIDJoypadRight] = 1
			}
			if pad.Axes[glfw.AxisLeftY] > 0.5 {
				state[p][lr.DeviceIDJoypadDown] = 1
			} else if pad.Axes[glfw.AxisLeftY] < -0.5 {
				state[p][lr.DeviceIDJoypadUp] = 1
			}
		}
		p++
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

	if device == lr.DeviceJoypad {
		if id >= uint(ActionLast) || index > 0 {
			return 0
		}
		return NewState[port][id]
	}
	if device == lr.DeviceAnalog {
		if index > uint(lr.DeviceIndexAnalogRight) || id > uint(lr.DeviceIDAnalogY) {
			return 0
		}

		return NewAnalogState[port][index][id]
	}

	return 0
}
