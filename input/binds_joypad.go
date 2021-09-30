package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/libretro"
)

var joyBinds = map[glfw.GamepadButton]uint32{
	glfw.ButtonDpadUp:    libretro.DeviceIDJoypadUp,
	glfw.ButtonDpadDown:  libretro.DeviceIDJoypadDown,
	glfw.ButtonDpadLeft:  libretro.DeviceIDJoypadLeft,
	glfw.ButtonDpadRight: libretro.DeviceIDJoypadRight,

	glfw.ButtonCircle:   libretro.DeviceIDJoypadA,
	glfw.ButtonCross:    libretro.DeviceIDJoypadB,
	glfw.ButtonSquare:   libretro.DeviceIDJoypadY,
	glfw.ButtonTriangle: libretro.DeviceIDJoypadX,

	glfw.ButtonLeftBumper:  libretro.DeviceIDJoypadL,
	glfw.ButtonRightBumper: libretro.DeviceIDJoypadR,

	glfw.ButtonLeftThumb:  libretro.DeviceIDJoypadL3,
	glfw.ButtonRightThumb: libretro.DeviceIDJoypadR3,

	glfw.ButtonStart: libretro.DeviceIDJoypadStart,
	glfw.ButtonBack:  libretro.DeviceIDJoypadSelect,
	glfw.ButtonGuide: ActionMenuToggle,
}
