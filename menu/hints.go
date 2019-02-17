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
