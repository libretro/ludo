// Package input exposes the two input callbacks Poll and State needed by the
// libretro implementation. It uses GLFW to access keyboard and joypad, and
// takes care of binding and auto configuring joypads.
package input

import (
	"log"

	deepcopy "github.com/barkimedes/go-deepcopy"
	"github.com/go-gl/glfw/v3.3/glfw"
	lr "github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

// MaxPlayers is the maximum number of players to poll input for
const MaxPlayers = 5

// MaxFrames is the max number of frames to keep in the input buffer. Used by netplay.
const MaxFrames = int64(60)

// LocalPlayerPort is the joypad port of the local player
var LocalPlayerPort = uint(0)

// RemotePlayerPort is the joypad port of the remote player
var RemotePlayerPort = uint(1)

var buffers = [MaxPlayers][MaxFrames]PlayerState{}

type joybinds map[bind]uint32

const btn = 0
const axis = 1

type bind struct {
	kind      uint32
	index     uint32
	direction float32
	threshold float32
}

// PlayerState is the state of inputs for a single player
type PlayerState [ActionLast]int16

// States can store the state of inputs for all players
type States [MaxPlayers]PlayerState

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

var oldMouseX float64
var oldMouseY float64

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

func index(offset int64) int64 {
	tick := state.Tick
	tick += offset
	return (MaxFrames + tick) % MaxFrames
}

// Serialize saves the current input state, used by netplay
func Serialize() [MaxPlayers][MaxFrames]PlayerState {
	copy := deepcopy.MustAnything(buffers)
	return copy.([MaxPlayers][MaxFrames]PlayerState)
}

// Unserialize restaures the input state from a save, used by netplay
func Unserialize(st interface{}) {
	copy := deepcopy.MustAnything(st)
	buffers = copy.([MaxPlayers][MaxFrames]PlayerState)
}

func getState(port uint, tick int64) PlayerState {
	frame := (MaxFrames + tick) % MaxFrames
	st := buffers[port][frame]
	return st
}

// GetLatest returns the most recent polled inputs
func GetLatest(port uint) PlayerState {
	return NewState[port]
}

func currentState(port uint) PlayerState {
	return getState(port, state.Tick)
}

// SetState forces the input state for a given player
func SetState(port uint, st PlayerState) {
	for i, b := range st {
		buffers[port][index(0)][i] = b
	}
}

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
// pollJoypads process joypads of all players
func pollJoypads() {
	p := LocalPlayerPort
	joy := glfw.Joystick(0)

	if !joy.IsGamepad() {
		return
	}
	pad := joy.GetGamepadState()
	if pad == nil {
		return
	}

	// mapping pad buttons
	for k, v := range joyBinds {
		if pad.Buttons[k] == glfw.Press {
			NewState[p][v] = 1
		}
	}

	// mapping pad triggers
	if pad.Axes[glfw.AxisLeftTrigger] > 0.5 {
		NewState[p][lr.DeviceIDJoypadL2] = 1
	}
	if pad.Axes[glfw.AxisRightTrigger] > 0.5 {
		NewState[p][lr.DeviceIDJoypadR2] = 1
	}

	// mapping analog sticks
	NewAnalogState[p][lr.DeviceIndexAnalogLeft][lr.DeviceIDAnalogX] = floatToAnalog(pad.Axes[glfw.AxisLeftX])
	NewAnalogState[p][lr.DeviceIndexAnalogLeft][lr.DeviceIDAnalogY] = floatToAnalog(pad.Axes[glfw.AxisLeftY])
	NewAnalogState[p][lr.DeviceIndexAnalogRight][lr.DeviceIDAnalogX] = floatToAnalog(pad.Axes[glfw.AxisRightX])
	NewAnalogState[p][lr.DeviceIndexAnalogRight][lr.DeviceIDAnalogY] = floatToAnalog(pad.Axes[glfw.AxisRightY])

	// optionally mapping analog sticks to dpad
	if settings.Current.MapAxisToDPad {
		if pad.Axes[glfw.AxisLeftX] < -0.5 {
			NewState[p][lr.DeviceIDJoypadLeft] = 1
		} else if pad.Axes[glfw.AxisLeftX] > 0.5 {
			NewState[p][lr.DeviceIDJoypadRight] = 1
		}
		if pad.Axes[glfw.AxisLeftY] > 0.5 {
			NewState[p][lr.DeviceIDJoypadDown] = 1
		} else if pad.Axes[glfw.AxisLeftY] < -0.5 {
			NewState[p][lr.DeviceIDJoypadUp] = 1
		}
	}
}

// pollKeyboard processes keyboard keys
func pollKeyboard() {
	for k, v := range keyBinds {
		if vid.Window.GetKey(k) == glfw.Press {
			NewState[LocalPlayerPort][v] = 1
		}
	}
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
	pollKeyboard()
	pollJoypads()

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

	// log.Println("input:", port, id, currentState(port)[id], state.Tick)
	if device == lr.DeviceJoypad {
		if id >= uint(ActionLast) || index > 0 {
			return 0
		}
		return currentState(port)[id]
	}
	if device == lr.DeviceAnalog {
		if index > uint(lr.DeviceIndexAnalogRight) || id > uint(lr.DeviceIDAnalogY) {
			return 0
		}

		return NewAnalogState[port][index][id]
	}

	if device == lr.DeviceMouse {
		x, y := vid.Window.GetCursorPos()
		if id == uint(lr.DeviceIDMouseX) {
			d := x - oldMouseX
			oldMouseX = x
			return int16(d)
		}
		if id == uint(lr.DeviceIDMouseY) {
			d := y - oldMouseY
			oldMouseY = y
			return int16(d)
		}
		if id == uint(lr.DeviceIDMouseLeft) && vid.Window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
			return 1
		}
		if id == uint(lr.DeviceIDMouseRight) && vid.Window.GetMouseButton(glfw.MouseButton2) == glfw.Press {
			return 1
		}
	}

	return 0
}
