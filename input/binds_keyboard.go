package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/libretro"
)

var keyBinds = map[glfw.Key]uint32{
	glfw.KeyX:          libretro.DeviceIDJoypadA,
	glfw.KeyZ:          libretro.DeviceIDJoypadB,
	glfw.KeyA:          libretro.DeviceIDJoypadY,
	glfw.KeyS:          libretro.DeviceIDJoypadX,
	glfw.KeyQ:          libretro.DeviceIDJoypadL,
	glfw.KeyW:          libretro.DeviceIDJoypadR,
	glfw.KeyUp:         libretro.DeviceIDJoypadUp,
	glfw.KeyDown:       libretro.DeviceIDJoypadDown,
	glfw.KeyLeft:       libretro.DeviceIDJoypadLeft,
	glfw.KeyRight:      libretro.DeviceIDJoypadRight,
	glfw.KeyEnter:      libretro.DeviceIDJoypadStart,
	glfw.KeyRightShift: libretro.DeviceIDJoypadSelect,
	glfw.KeySpace:      ActionFastForwardToggle,
	glfw.KeyP:          ActionMenuToggle,
	glfw.KeyF:          ActionFullscreenToggle,
	glfw.KeyEscape:     ActionShouldClose,
}
