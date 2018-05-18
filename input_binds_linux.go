package main

import "libretro"

var joyBinds = map[string]joybinds{
	"Xbox 360 Wired Controller": xbox360JoyBinds,
}

// Joypad bindings fox Xbox360 pad on Linux
var xbox360JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:      libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:      libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:      libretro.DeviceIDJoypadY,
	bind{btn, 3, 0, 0}:      libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:      libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:      libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:      libretro.DeviceIDJoypadSelect,
	bind{btn, 7, 0, 0}:      libretro.DeviceIDJoypadStart,
	bind{btn, 8, 0, 0}:      menuActionMenuToggle,
	bind{btn, 9, 0, 0}:      libretro.DeviceIDJoypadL3,
	bind{btn, 10, 0, 0}:     libretro.DeviceIDJoypadR3,
	bind{axis, 6, -1, -0.5}: libretro.DeviceIDJoypadLeft,
	bind{axis, 6, 1, 0.5}:   libretro.DeviceIDJoypadRight,
	bind{axis, 7, -1, -0.5}: libretro.DeviceIDJoypadUp,
	bind{axis, 7, 1, 0.5}:   libretro.DeviceIDJoypadDown,
	bind{axis, 2, 1, 0.5}:   libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0.5}:   libretro.DeviceIDJoypadR2,
}
