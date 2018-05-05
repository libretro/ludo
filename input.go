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

var buttonBinds = map[byte]uint32{
	0:  libretro.DeviceIDJoypadUp,
	1:  libretro.DeviceIDJoypadDown,
	2:  libretro.DeviceIDJoypadLeft,
	3:  libretro.DeviceIDJoypadRight,
	4:  libretro.DeviceIDJoypadStart,
	5:  libretro.DeviceIDJoypadSelect,
	6:  libretro.DeviceIDJoypadL3,
	7:  libretro.DeviceIDJoypadR3,
	8:  libretro.DeviceIDJoypadL,
	9:  libretro.DeviceIDJoypadR,
	10: menuActionMenuToggle,
	11: libretro.DeviceIDJoypadB,
	12: libretro.DeviceIDJoypadA,
	13: libretro.DeviceIDJoypadY,
	14: libretro.DeviceIDJoypadX,
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
		if len(buttonState) > 0 {
			for k, v := range buttonBinds {
				// Ensure the button is available in one of the states.
				if int(k) < len(buttonState) {
					if glfw.Action(buttonState[k]) == glfw.Press {
						newState[p][v] = true
					}
				} else {
					fmt.Println("Unknown button index: ", k)
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
