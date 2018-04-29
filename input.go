package main

import (
	"C"

	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
#include "libretro.h"
*/
import "C"

var binds = map[glfw.Key]C.int{
	glfw.KeyX:         C.RETRO_DEVICE_ID_JOYPAD_A,
	glfw.KeyZ:         C.RETRO_DEVICE_ID_JOYPAD_B,
	glfw.KeyA:         C.RETRO_DEVICE_ID_JOYPAD_Y,
	glfw.KeyS:         C.RETRO_DEVICE_ID_JOYPAD_X,
	glfw.KeyUp:        C.RETRO_DEVICE_ID_JOYPAD_UP,
	glfw.KeyDown:      C.RETRO_DEVICE_ID_JOYPAD_DOWN,
	glfw.KeyLeft:      C.RETRO_DEVICE_ID_JOYPAD_LEFT,
	glfw.KeyRight:     C.RETRO_DEVICE_ID_JOYPAD_RIGHT,
	glfw.KeyEnter:     C.RETRO_DEVICE_ID_JOYPAD_START,
	glfw.KeyBackspace: C.RETRO_DEVICE_ID_JOYPAD_SELECT,
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
