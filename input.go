package main

import (
	"C"

	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
#include "libretro.h"
*/
import "C"

var binds = map[glfw.Key]uint32{
	glfw.KeyX:         retroDeviceIDJoypadA,
	glfw.KeyZ:         retroDeviceIDJoypadB,
	glfw.KeyA:         retroDeviceIDJoypadY,
	glfw.KeyS:         retroDeviceIDJoypadX,
	glfw.KeyUp:        retroDeviceIDJoypadUp,
	glfw.KeyDown:      retroDeviceIDJoypadDown,
	glfw.KeyLeft:      retroDeviceIDJoypadLeft,
	glfw.KeyRight:     retroDeviceIDJoypadRight,
	glfw.KeyEnter:     retroDeviceIDJoypadStart,
	glfw.KeyBackspace: retroDeviceIDJoypadSelect,
}

var joy [C.RETRO_DEVICE_ID_JOYPAD_R3 + 1]bool

func inputPoll() {
	for k, v := range binds {
		joy[v] = (window.GetKey(k) == glfw.Press)
	}

	// Close the window when the user hits the Escape key.
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func inputState(port C.unsigned, device C.unsigned, index C.unsigned, id C.unsigned) C.int16_t {
	if port > 0 || index > 0 || device != C.RETRO_DEVICE_JOYPAD {
		return 0
	}

	if id < 255 && joy[id] {
		return 1
	}
	return 0
}
