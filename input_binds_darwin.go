package main

import "libretro"

// Joypad bindings fox Xbox360 pad on OSX
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
