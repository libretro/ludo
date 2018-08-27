package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/notifications"
	"github.com/libretro/go-playthemall/video"
)

const numPlayers = 5

type joybinds map[bind]uint32

var keyBinds = map[glfw.Key]uint32{
	glfw.KeyX:         libretro.DeviceIDJoypadA,
	glfw.KeyZ:         libretro.DeviceIDJoypadB,
	glfw.KeyA:         libretro.DeviceIDJoypadY,
	glfw.KeyS:         libretro.DeviceIDJoypadX,
	glfw.KeyUp:        libretro.DeviceIDJoypadUp,
	glfw.KeyDown:      libretro.DeviceIDJoypadDown,
	glfw.KeyLeft:      libretro.DeviceIDJoypadLeft,
	glfw.KeyRight:     libretro.DeviceIDJoypadRight,
	glfw.KeyEnter:     libretro.DeviceIDJoypadStart,
	glfw.KeyBackspace: libretro.DeviceIDJoypadSelect,
	glfw.KeyP:         ActionMenuToggle,
	glfw.KeyF:         ActionFullscreenToggle,
	glfw.KeyEscape:    ActionShouldClose,
}

const btn = 0
const axis = 1

type bind struct {
	kind      uint32
	index     uint32
	direction float32
	threshold float32
}

type inputstate [numPlayers][ActionLast]bool

// Input state for all the players
var (
	NewState inputstate // input state for the current frame
	OldState inputstate // input state for the previous frame
	Released inputstate // keys just released during this frame
	Pressed  inputstate // keys just pressed during this frame
)

// joystickCallback is triggered when a joypad is plugged.
func joystickCallback(joy int, event int) {
	switch event {
	case 262145:
		notifications.DisplayAndLog("Input", "Joystick #%d plugged: %s.", joy, glfw.GetJoystickName(glfw.Joystick(joy)))
	case 262146:
		notifications.DisplayAndLog("Input", "Joystick #%d unplugged.", joy)
	default:
		notifications.DisplayAndLog("Input", "Joystick #%d unhandled event: %d.", joy, event)
	}
}

// ContextReseter is an interface to to allow reloading icons after the
// window is recreated when switching fullscreen
type ContextReseter interface {
	ContextReset()
}

var vid *video.Video
var menu ContextReseter

// Init initializes the input package
func Init(v *video.Video, m ContextReseter) {
	vid = v
	menu = m
	glfw.SetJoystickCallback(joystickCallback)
}

// Resets all retropad buttons to false
func reset(state inputstate) inputstate {
	for p := range state {
		for k := range state[p] {
			state[p][k] = false
		}
	}
	return state
}

// pollJoypads process joypads of all players
func pollJoypads(state inputstate) inputstate {
	for p := range state {
		buttonState := glfw.GetJoystickButtons(glfw.Joystick(p))
		axisState := glfw.GetJoystickAxes(glfw.Joystick(p))
		name := glfw.GetJoystickName(glfw.Joystick(p))
		jb := joyBinds[name]
		if len(buttonState) > 0 {
			for k, v := range jb {
				switch k.kind {
				case btn:
					if glfw.Action(buttonState[k.index]) == glfw.Press {
						state[p][v] = true
					}
				case axis:
					if k.direction*axisState[k.index] > k.threshold*k.direction {
						state[p][v] = true
					}
				}
			}
		}
	}
	return state
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
	NewState = reset(NewState)
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
