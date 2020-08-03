package menu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/video"
)

// Used to easily compose different hint bars based on the context.
func stackHintLeft(stack *float32, icon uint32, label string, h int) {
	c := video.Color{R: 0.28, G: 0.28, B: 0.28, A: 1}
	vid.Font.SetColor(0.28, 0.28, 0.28, 1.0)
	vid.DrawImage(icon, *stack, float32(h)-79*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, 0, c)
	*stack += 70 * menu.ratio
	vid.Font.Printf(*stack, float32(h)-30*menu.ratio, 0.5*menu.ratio, label)
	*stack += vid.Font.Width(0.5*menu.ratio, label)
	*stack += 32 * menu.ratio
}

// Used to easily compose different hint bars based on the context.
func stackHintRight(stack *float32, icon uint32, label string, h int) {
	c := video.Color{R: 0.28, G: 0.28, B: 0.28, A: 1}
	*stack -= vid.Font.Width(0.5*menu.ratio, label)
	vid.Font.SetColor(0.28, 0.28, 0.28, 1.0)
	vid.Font.Printf(*stack, float32(h)-30*menu.ratio, 0.5*menu.ratio, label)
	*stack -= 70 * menu.ratio
	vid.DrawImage(icon, *stack, float32(h)-79*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, 0, c)
	*stack -= 32 * menu.ratio
}

func hintIcons() (arrows, upDown, leftRight, a, b, x, y, start, slct, guide uint32) {
	if glfw.Joystick(0).Present() && input.HasBinding(glfw.Joystick(0)) {
		return menu.icons["pad-arrows"],
			menu.icons["pad-up-down"],
			menu.icons["pad-left-right"],
			menu.icons["pad-a"],
			menu.icons["pad-b"],
			menu.icons["pad-x"],
			menu.icons["pad-y"],
			menu.icons["pad-start"],
			menu.icons["pad-select"],
			menu.icons["pad-guide"]
	}
	return menu.icons["key-arrows"],
		menu.icons["key-up-down"],
		menu.icons["key-left-right"],
		menu.icons["key-x"],
		menu.icons["key-z"],
		menu.icons["key-s"],
		menu.icons["key-a"],
		menu.icons["key-return"],
		menu.icons["key-shift"],
		menu.icons["key-p"]
}
