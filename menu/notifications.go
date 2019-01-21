package menu

import (
	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/video"
)

// RenderNotifications draws the list of notification messages on the viewport
func RenderNotifications() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)
	var h float32 = 75
	stack := h
	for _, n := range notifications.List() {
		fading := float32(n.Frames) / 120
		if fading > 1 {
			fading = 1
		}
		offset := fading*20 - 20
		lw := vid.Font.Width(0.5*menu.ratio, n.Message)
		vid.DrawRoundedRect(
			25*menu.ratio,
			(stack+offset-46)*menu.ratio,
			lw+40*menu.ratio,
			70*menu.ratio,
			0.25,
			video.Color{R: 0.4, G: 0.4, B: 0, A: fading},
		)
		vid.Font.SetColor(1, 1, 0.85, fading)
		vid.Font.Printf(
			45*menu.ratio,
			(stack+offset)*menu.ratio,
			0.5*menu.ratio,
			n.Message,
		)
		stack += h + offset
	}
}
