package main

import (
	"fmt"
	"libretro"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const numPlayers = 5

const (
	menuActionMenuToggle       uint32 = libretro.DeviceIDJoypadR3 + 1
	menuActionFullscreenToggle uint32 = libretro.DeviceIDJoypadR3 + 2
	menuActionShouldClose      uint32 = libretro.DeviceIDJoypadR3 + 3
	menuActionLast             uint32 = libretro.DeviceIDJoypadR3 + 4
)

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
	glfw.KeyP:         menuActionMenuToggle,
	glfw.KeyF:         menuActionFullscreenToggle,
	glfw.KeyEscape:    menuActionShouldClose,
}

const btn = 0
const axis = 1

type bind struct {
	kind      uint32
	index     uint32
	direction float32
	threshold float32
}

// Joypad binding fox Xbox360 pad on OSX
var joyBinds = map[bind]uint32{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadUp,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadDown,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadLeft,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadRight,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadL3,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadR3,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 10, 0, 0}: menuActionMenuToggle,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadB,
	bind{btn, 12, 0, 0}: libretro.DeviceIDJoypadA,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadY,
	bind{btn, 14, 0, 0}: libretro.DeviceIDJoypadX,
	bind{axis, 4, 1, 0}: libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0}: libretro.DeviceIDJoypadR2,
	// Uncomment this to bind left analog to directions
	// bind{axis, 0, -1, -0.5}: libretro.DeviceIDJoypadLeft,
	// bind{axis, 0, 1, 0.5}:   libretro.DeviceIDJoypadRight,
	// bind{axis, 1, -1, -0.5}: libretro.DeviceIDJoypadUp,
	// bind{axis, 1, 1, 0.5}:   libretro.DeviceIDJoypadDown,
}

// Input state for all the players
var (
	newState [numPlayers][menuActionLast]bool // input state for the current frame
	oldState [numPlayers][menuActionLast]bool // input state for the previous frame
	released [numPlayers][menuActionLast]bool // keys just released during this frame
	pressed  [numPlayers][menuActionLast]bool // keys just pressed during this frame
)

func joystickCallback(joy int, event int) {
	var message string
	switch event {
	case 262145:
		message = fmt.Sprintf("Joystick #%d plugged: %s.", joy, glfw.GetJoystickName(glfw.Joystick(joy)))
		break
	case 262146:
		message = fmt.Sprintf("Joystick #%d unplugged.", joy)
		break
	default:
		message = fmt.Sprintf("Joystick #%d unhandled event: %d.", joy, event)
	}
	fmt.Printf("[Input]: %s\n", message)
	notify(message, 240)
}

func inputInit() {
	glfw.SetJoystickCallback(joystickCallback)
}

func inputPoll() {
	// Reset all retropad buttons to false
	for p := range newState {
		for k := range newState[p] {
			newState[p][k] = false
		}
	}

	// Process joypads of all players
	for p := range newState {
		buttonState := glfw.GetJoystickButtons(glfw.Joystick(p))
		axisState := glfw.GetJoystickAxes(glfw.Joystick(p))
		if len(buttonState) > 0 {
			for k, v := range joyBinds {
				switch k.kind {
				case btn:
					if glfw.Action(buttonState[k.index]) == glfw.Press {
						newState[p][v] = true
					}
				case axis:
					if k.direction*axisState[k.index] > k.threshold*k.direction {
						newState[p][v] = true
					}
				}
			}
		}
	}

	// Process keyboard keys
	for k, v := range keyBinds {
		if window.GetKey(k) == glfw.Press {
			newState[0][v] = true
		}
	}

	// Compute the keys pressed or released during this frame
	for p := range newState {
		for k := range newState[p] {
			pressed[p][k] = newState[p][k] && !oldState[p][k]
			released[p][k] = !newState[p][k] && oldState[p][k]
		}
	}

	// Toggle the menu if menuActionMenuToggle is pressed
	if released[0][menuActionMenuToggle] {
		g.menuActive = !g.menuActive
	}

	// Toggle fullscreen if menuActionFullscreenToggle is pressed
	if released[0][menuActionFullscreenToggle] {
		toggleFullscreen()
	}

	// Close on escape
	if pressed[0][menuActionShouldClose] {
		window.SetShouldClose(true)
	}

	// Store the old input state for comparisions
	oldState = newState
}

func inputState(port uint, device uint32, index uint, id uint) int16 {
	if id >= 255 || index > 0 || device != libretro.DeviceJoypad {
		return 0
	}

	if newState[port][id] {
		return 1
	}
	return 0
}
