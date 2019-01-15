package menu

import (
	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/video"
)

// RenderNotifications draws the list of notification messages on the viewport
func RenderNotifications() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)
	for i, n := range notifications.List() {
		a := float32(n.Frames) / 120
		lw := vid.Font.Width(0.5*menu.ratio, n.Message)
		vid.DrawRoundedRect(
			25*menu.ratio, (float32(75+75*i)-46)*menu.ratio,
			lw+40*menu.ratio, 70*menu.ratio, 0.25,
			video.Color{R: 0.4, G: 0.4, B: 0, A: a})
		vid.Font.SetColor(1, 1, 0.85, a)
		vid.Font.Printf(
			45*menu.ratio, float32(75+75*i)*menu.ratio, 0.5*menu.ratio, n.Message)
	}
}
