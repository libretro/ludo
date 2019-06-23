package menu

import "github.com/libretro/ludo/video"

// HintBar is the bar showing at the bottom of the screen
func HintBar(props *Props, children ...func()) func() {
	w, h := vid.Window.GetFramebufferSize()
	return HBox(&Props{
		Y:      float32(h) - 70*menu.ratio,
		Width:  float32(w),
		Height: 70 * menu.ratio,
		Color:  video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1},
		Hidden: props.Hidden,
	},
		children...,
	)
}

// Hint is a widget combining an icon and a label
func Hint(props *Props, icon string, title string) func() {
	darkGrey := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	return HBox(props,
		Box(&Props{Width: 15}),
		Image(&Props{
			Width:  70 * menu.ratio,
			Height: 70 * menu.ratio,
			Scale:  1,
			Color:  darkGrey,
		}, menu.icons[icon]),
		Label(&Props{
			Height: 70 * menu.ratio,
			Scale:  0.5 * menu.ratio,
			Color:  darkGrey,
		}, title),
		Box(&Props{Width: 15}),
	)
}
