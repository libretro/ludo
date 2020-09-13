// Package input exposes the two input callbacks Poll and State needed by the
// libretro implementation. It uses GLFW to access keyboard and joypad, and
// takes care of binding and auto configuring joypads.
package input

import (
	"reflect"
	"runtime"
	"strings"

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

type inputstate [MaxPlayers][ActionLast]bool

// Input state for all the players
var (
	NewState inputstate // input state for the current frame
	OldState inputstate // input state for the previous frame
	Released inputstate // keys just released during this frame
	Pressed  inputstate // keys just pressed during this frame
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

// pollJoypads process joypads of all players
func pollJoypads(state inputstate) inputstate {
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
						state[p][v] = true
					}
				case axis:
					idx := k.index
					idx = getJoystickIdx(&jb, name, idx)
					if jb[k] == libretro.DeviceIDJoypadR3 {
						idx += 2
					}
					processAxisState(idx, &axisState, &k, &state, p, v)

				}
			}
		}
	}
	return state
}

func processAxisState(idx uint32, axisState *[]float32, k *bind, state *inputstate, p int, v uint32) {
	if int(idx+1) < len(*axisState) &&
		((k.direction*(*axisState)[idx] > k.threshold*k.direction ||
			k.direction*(*axisState)[idx] < -k.threshold*k.direction) ||
			(k.direction*(*axisState)[idx+1] > k.threshold*k.direction ||
				k.direction*(*axisState)[idx+1] < -k.threshold*k.direction)) {
		state[p][v] = true
		if !settings.Current.MapAxisToDPad {
			return
		}
		switch {
		// idx == 0 for L3, idx == 3 for R3
		case idx == 0 && (*axisState)[idx] < -k.threshold:
			state[p][libretro.DeviceIDJoypadLeft] = true
		case idx == 0 && (*axisState)[idx] > k.threshold:
			state[p][libretro.DeviceIDJoypadRight] = true
		}
		switch {
		case idx == 0 && (*axisState)[idx+1] > k.threshold:
			state[p][libretro.DeviceIDJoypadDown] = true
		case idx == 0 && (*axisState)[idx+1] < -k.threshold:
			state[p][libretro.DeviceIDJoypadUp] = true
		}
	}
}

// getJoystickIndx returns X axis index for stick `idx` of joystick `name` in axisState
func getJoystickIdx(jb *joybinds, name string, idx uint32) uint32 {
	if strings.HasPrefix(name, "Sony") {
		if runtime.GOOS == "windows" {
			idx -= 10
		} else if runtime.GOOS == "linux" {
			idx -= 11
		} else if runtime.GOOS == "darwin" {
			switch reflect.TypeOf(jb) {
			case reflect.TypeOf(ds4JoyBinds):
				idx -= 10
			case reflect.TypeOf(ds3JoyBinds):
				idx--
			}
		}
	}
	return idx
}

// pollKeyboard processes keyboard keys
func pollKeyboard(state inputstate) inputstate {
	for k, v := range keyBinds {
		if vid.Window.GetKey(k) == glfw.Press {
			state[0][v] = true
		}
	}
	return state
}

// Compute the keys pressed or released during this frame
func getPressedReleased(new inputstate, old inputstate) (inputstate, inputstate) {
	for p := range new {
		for k := range new[p] {
			Pressed[p][k] = new[p][k] && !old[p][k]
			Released[p][k] = !new[p][k] && old[p][k]
		}
	}
	return Pressed, Released
}

// Poll calculates the input state. It is meant to be called for each frame.
func Poll() {
	NewState = inputstate{}
	NewState = pollJoypads(NewState)
	NewState = pollKeyboard(NewState)
	Pressed, Released = getPressedReleased(NewState, OldState)

	// Store the old input state for comparisions
	OldState = NewState
}

// State is a callback passed to core.SetInputState
// It returns 1 if the button corresponding to the parameters is pressed
func State(port uint, device uint32, index uint, id uint) int16 {
	if id >= 255 || index > 0 || device != libretro.DeviceJoypad {
		return 0
	}

	if NewState[port][id] {
		return 1
	}
	return 0
}

// HasBinding returns true if the joystick has an autoconfig binding
func HasBinding(joy glfw.Joystick) bool {
	name := glfw.Joystick.GetName(joy)
	_, ok := joyBinds[name]
	return ok
}
