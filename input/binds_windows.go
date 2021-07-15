package input

import "github.com/libretro/ludo/libretro"

var joyBinds = map[string]joybinds{
	"Microsoft X-Box 360 pad":    xbox360JoyBinds,
	"Xbox 360 Controller":        xboxOneJoyBinds,
	"Xbox Controller":            xboxOneJoyBinds,
	"Wireless Xbox Controller":   xboxOneJoyBinds,
	"Wireless Controller":        ds4JoyBinds,
	"PLAYSTATION(R)3 Controller": ds3JoyBinds,
	"8Bitdo NES30 Pro":           nes30proJoyBinds,
	"SFC30 Joystick":             sfc30JoyBinds,
}

var xbox360JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadL3,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadR3,
	bind{btn, 10, 0, 0}: libretro.DeviceIDJoypadUp,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadRight,
	bind{btn, 12, 0, 0}: libretro.DeviceIDJoypadDown,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadLeft,
	bind{axis, 4, 1, 0}: libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0}: libretro.DeviceIDJoypadL3,
}

var xboxOneJoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadL3,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadR3,
	bind{btn, 10, 0, 0}: libretro.DeviceIDJoypadUp,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadRight,
	bind{btn, 12, 0, 0}: libretro.DeviceIDJoypadDown,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadLeft,
	bind{axis, 4, 1, 0}: libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0}: libretro.DeviceIDJoypadL3,
}

var ds4JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadL2,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadR2,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 10, 0, 0}: libretro.DeviceIDJoypadL3,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadR3,
	bind{btn, 12, 0, 0}: ActionMenuToggle,
	bind{btn, 14, 0, 0}: libretro.DeviceIDJoypadUp,
	bind{btn, 15, 0, 0}: libretro.DeviceIDJoypadRight,
	bind{btn, 16, 0, 0}: libretro.DeviceIDJoypadDown,
	bind{btn, 17, 0, 0}: libretro.DeviceIDJoypadLeft,
}

// Detected but doesn't send inputs
var ds3JoyBinds = joybinds{}

// Joypad bindings fox the 8BITDO NES30 PRO GamePad (Wired) on Windows
var nes30proJoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadL2,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadR2,
	bind{btn, 10, 0, 0}: libretro.DeviceIDJoypadSelect,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadStart,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadL3,
	bind{btn, 14, 0, 0}: libretro.DeviceIDJoypadR3,
	bind{btn, 17, 0, 0}: libretro.DeviceIDJoypadDown,
	bind{btn, 18, 0, 0}: libretro.DeviceIDJoypadLeft,
	bind{btn, 16, 0, 0}: libretro.DeviceIDJoypadRight,
	bind{btn, 15, 0, 0}: libretro.DeviceIDJoypadUp,
}

// Joypad bindings fox the 8BITDO SFC30 pad (Wired) on Windows
var sfc30JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:      libretro.DeviceIDJoypadA,
	bind{btn, 1, 0, 0}:      libretro.DeviceIDJoypadB,
	bind{btn, 3, 0, 0}:      libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:      libretro.DeviceIDJoypadY,
	bind{btn, 6, 0, 0}:      libretro.DeviceIDJoypadL,
	bind{btn, 7, 0, 0}:      libretro.DeviceIDJoypadR,
	bind{btn, 10, 0, 0}:     libretro.DeviceIDJoypadSelect,
	bind{btn, 11, 0, 0}:     libretro.DeviceIDJoypadStart,
	bind{axis, 0, -1, -0.5}: libretro.DeviceIDJoypadLeft,
	bind{axis, 1, -1, -0.5}: libretro.DeviceIDJoypadUp,
	bind{axis, 0, 1, 0.5}:   libretro.DeviceIDJoypadRight,
	bind{axis, 1, 1, 0.5}:   libretro.DeviceIDJoypadDown,
}
