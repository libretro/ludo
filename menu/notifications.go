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
		lw := vid.Font.Width(0.6*menu.ratio, n.Message)
		vid.DrawRoundedRect(
			25*menu.ratio, (float32(80+80*i)-52)*menu.ratio,
			lw+35*menu.ratio, 75*menu.ratio, 0.25,
			video.Color{R: 0.4, G: 0.4, B: 0, A: float32(n.Frames) / 120.0})
		vid.Font.SetColor(1, 1, 0.85, float32(n.Frames)/120.0)
		vid.Font.Printf(
			45*menu.ratio, float32(80+80*i)*menu.ratio, 0.6*menu.ratio, n.Message)
	}
}
