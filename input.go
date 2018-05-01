package main

import (
	"fmt"
	"libretro"

	"github.com/go-gl/glfw/v3.2/glfw"
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
}

var buttonBinds = map[byte]uint32{
	0: libretro.DeviceIDJoypadUp,
	1: libretro.DeviceIDJoypadDown,
	2: libretro.DeviceIDJoypadLeft,
	3: libretro.DeviceIDJoypadRight,
	4: libretro.DeviceIDJoypadStart,
	5: libretro.DeviceIDJoypadSelect,
	6: libretro.DeviceIDJoypadL3,
	7: libretro.DeviceIDJoypadR3,
	8: libretro.DeviceIDJoypadL,
	9: libretro.DeviceIDJoypadR,
	//10: libretro.DeviceIDJoypadGuide,
	11: libretro.DeviceIDJoypadB,
	12: libretro.DeviceIDJoypadA,
	13: libretro.DeviceIDJoypadY,
	14: libretro.DeviceIDJoypadX,
}

var retroPad [libretro.DeviceIDJoypadR3 + 1]bool

func joystickCallback(joy int, event int) {
	switch event {
	case 262145:
		fmt.Printf("[Input]: Joystick #%d plugged.\n", joy)
		break
	case 262146:
		fmt.Printf("[Input]: Joystick #%d unplugged.\n", joy)
		break
	default:
		fmt.Printf("[Input]: Joystick #%d unhandled event: %d.\n", joy, event)
	}
}

func inputInit() {
	glfw.SetJoystickCallback(joystickCallback)
}

func inputPoll() {
	for i := range retroPad {
		retroPad[i] = false
	}

	buttonState := glfw.GetJoystickButtons(glfw.Joystick1)

	if len(buttonState) > 0 {
		for k, v := range buttonBinds {
			if glfw.Action(buttonState[k]) == glfw.Press {
				retroPad[v] = true
			}
		}
	}

	for k, v := range keyBinds {
		if window.GetKey(k) == glfw.Press {
			retroPad[v] = true
		}
	}

	// Close the window when the user hits the Escape key.
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func inputState(port uint, device uint32, index uint, id uint) int16 {
	if port > 0 || index > 0 || device != libretro.DeviceJoypad {
		return 0
	}

	if id < 255 && retroPad[id] {
		return 1
	}
	return 0
}
