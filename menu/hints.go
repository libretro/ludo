package menu

import "github.com/libretro/ludo/video"

// Used to easily compose different hint bars based on the context.
func stackHint(stack *float32, icon, label string, h int) {
	c := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	vid.Font.SetColor(0.25, 0.25, 0.25, 1.0)
	*stack += 30 * menu.ratio
	vid.DrawImage(menu.icons[icon], *stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	*stack += 70 * menu.ratio
	vid.Font.Printf(*stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, label)
	*stack += vid.Font.Width(0.5*menu.ratio, label)
}

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
