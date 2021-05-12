package menu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/input"
)

// HintBar is the bar showing at the bottom of the screen
func HintBar(props *Props, children ...func()) func() {
	w, h := vid.Window.GetFramebufferSize()
	return HBox(&Props{
		Y:      float32(h) - 70*menu.ratio,
		Width:  float32(w),
		Height: 70 * menu.ratio,
		Color:  lightGrey,
		Hidden: props.Hidden,
	},
		children...,
	)
}

// Hint is a widget combining an icon and a label
func Hint(props *Props, icon uint32, title string) func() {
	darkGrey := darkGrey
	return HBox(props,
		Box(&Props{Width: 15}),
		Image(&Props{
			Width:  70 * menu.ratio,
			Height: 70 * menu.ratio,
			Scale:  1,
			Color:  darkGrey,
		}, icon),
		Label(&Props{
			Height: 70 * menu.ratio,
			Scale:  0.4 * menu.ratio,
			Color:  darkGrey,
		}, title),
		Box(&Props{Width: 15}),
	)
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
