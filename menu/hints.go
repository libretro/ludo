package menu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Used to easily compose different hint bars based on the context.
func stackHint(stack *float32, icon uint32, label string, h int) {
	menu.Font.SetColor(darkGrey)
	*stack += 30 * menu.ratio
	menu.DrawImage(icon, *stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, darkGrey)
	*stack += 70 * menu.ratio
	menu.Font.Printf(*stack, float32(h)-23*menu.ratio, 0.4*menu.ratio, label)
	*stack += menu.Font.Width(0.4*menu.ratio, label)
}

func hintIcons() (arrows, upDown, leftRight, a, b, x, y, start, slct, guide uint32) {
	if glfw.Joystick(0).IsGamepad() {
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
