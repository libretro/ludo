package input

import "github.com/libretro/ludo/libretro"

var joyBinds = map[string]joybinds{
	"Microsoft X-Box 360 pad":                            xbox360JoyBinds,
	"Microsoft X-Box One S pad":                          xboxOneJoyBinds,
	"Sony Computer Entertainment Wireless Controller":    ds4JoyBinds,
	"Sony Interactive Entertainment Wireless Controller": ds4JoyBinds,
	"Wireless Controller":                                ds4JoyBinds,
	"Sony PLAYSTATION(R)3 Controller":                    ds3JoyBinds,
	"SHANWAN PS3 GamePad":                                ds3JoyBinds,
	"8Bitdo NES30 Pro   8Bitdo NES30 Pro":                nes30proJoyBinds,
	"SFC30              SFC30 Joystick":                  sfc30JoyBinds,
	"SNES30             SNES30 Joy    ":                  sfc30JoyBinds,
}

// Joypad bindings fox Xbox360 pad on Linux
var xbox360JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:    libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:    libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:    libretro.DeviceIDJoypadY,
	bind{btn, 3, 0, 0}:    libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:    libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:    libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:    libretro.DeviceIDJoypadSelect,
	bind{btn, 7, 0, 0}:    libretro.DeviceIDJoypadStart,
	bind{btn, 8, 0, 0}:    ActionMenuToggle,
	bind{btn, 9, 0, 0}:    libretro.DeviceIDJoypadL3,
	bind{btn, 10, 0, 0}:   libretro.DeviceIDJoypadR3,
	bind{btn, 11, 1, 0}:   libretro.DeviceIDJoypadUp,
	bind{btn, 12, 1, 0}:   libretro.DeviceIDJoypadRight,
	bind{btn, 13, 0, 0}:   libretro.DeviceIDJoypadDown,
	bind{btn, 14, 0, 0}:   libretro.DeviceIDJoypadLeft,
	bind{axis, 2, 1, 0.5}: libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0.5}: libretro.DeviceIDJoypadR2,
}

// Joypad bindings fox Xbox One pad on Linux
var xboxOneJoyBinds = joybinds{
	bind{btn, 0, 0, 0}:    libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:    libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:    libretro.DeviceIDJoypadY,
	bind{btn, 3, 0, 0}:    libretro.DeviceIDJoypadX,
	bind{btn, 4, 0, 0}:    libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:    libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:    libretro.DeviceIDJoypadSelect,
	bind{btn, 7, 0, 0}:    libretro.DeviceIDJoypadStart,
	bind{btn, 8, 0, 0}:    ActionMenuToggle,
	bind{btn, 9, 0, 0}:    libretro.DeviceIDJoypadL3,
	bind{btn, 10, 0, 0}:   libretro.DeviceIDJoypadR3,
	bind{btn, 11, 1, 0}:   libretro.DeviceIDJoypadUp,
	bind{btn, 12, 1, 0}:   libretro.DeviceIDJoypadRight,
	bind{btn, 13, 0, 0}:   libretro.DeviceIDJoypadDown,
	bind{btn, 14, 0, 0}:   libretro.DeviceIDJoypadLeft,
	bind{axis, 2, 1, 0.5}: libretro.DeviceIDJoypadL2,
	bind{axis, 5, 1, 0.5}: libretro.DeviceIDJoypadR2,
}

// Joypad bindings fox DualShock 4 pad on Linux
var ds4JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadL2,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadR2,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 10, 0, 0}: ActionMenuToggle,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadL3,
	bind{btn, 12, 0, 0}: libretro.DeviceIDJoypadR3,
	bind{btn, 16, 0, 0}: libretro.DeviceIDJoypadLeft,
	bind{btn, 14, 0, 0}: libretro.DeviceIDJoypadRight,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadUp,
	bind{btn, 15, 0, 0}: libretro.DeviceIDJoypadDown,
}

// Joypad bindings fox DualShock 3 pad on Linux
var ds3JoyBinds = joybinds{
	bind{btn, 0, 0, 0}:  libretro.DeviceIDJoypadB,
	bind{btn, 1, 0, 0}:  libretro.DeviceIDJoypadA,
	bind{btn, 2, 0, 0}:  libretro.DeviceIDJoypadX,
	bind{btn, 3, 0, 0}:  libretro.DeviceIDJoypadY,
	bind{btn, 4, 0, 0}:  libretro.DeviceIDJoypadL,
	bind{btn, 5, 0, 0}:  libretro.DeviceIDJoypadR,
	bind{btn, 6, 0, 0}:  libretro.DeviceIDJoypadL2,
	bind{btn, 7, 0, 0}:  libretro.DeviceIDJoypadR2,
	bind{btn, 8, 0, 0}:  libretro.DeviceIDJoypadSelect,
	bind{btn, 9, 0, 0}:  libretro.DeviceIDJoypadStart,
	bind{btn, 10, 0, 0}: ActionMenuToggle,
	bind{btn, 11, 0, 0}: libretro.DeviceIDJoypadL3,
	bind{btn, 12, 0, 0}: libretro.DeviceIDJoypadR3,
	bind{btn, 13, 0, 0}: libretro.DeviceIDJoypadUp,
	bind{btn, 14, 0, 0}: libretro.DeviceIDJoypadDown,
	bind{btn, 15, 0, 0}: libretro.DeviceIDJoypadLeft,
	bind{btn, 16, 0, 0}: libretro.DeviceIDJoypadRight,
}

// Joypad bindings fox the 8BITDO NES30 PRO GamePad (Wired) on Linux
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

// Joypad bindings fox the 8BITDO SFC30 pad (Wired) on Linux
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
