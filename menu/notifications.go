package menu

import (
	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/video"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var severityFgColor = map[string]colorful.Color{
	"error":   colorful.Hcl(40, 0.75, 0.85),
	"warning": colorful.Hcl(90, 0.75, 0.85),
	"success": colorful.Hcl(120, 0.75, 0.85),
	"info":    colorful.Hcl(230, 0.75, 0.85),
}

var severityBgColor = map[string]colorful.Color{
	"error":   colorful.Hcl(40, 0.65, 0.1),
	"warning": colorful.Hcl(90, 0.65, 0.1),
	"success": colorful.Hcl(120, 0.65, 0.1),
	"info":    colorful.Hcl(230, 0.65, 0.1),
}

// RenderNotifications draws the list of notification messages on the viewport
func RenderNotifications() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)
	var h float32 = 75
	stack := h
	for _, n := range notifications.List() {
		fading := float32(n.Frames) / 10
		if fading > 1 {
			fading = 1
		}
		offset := fading*h - h
		lw := vid.Font.Width(0.5*menu.ratio, n.Message)
		fg := severityFgColor[n.Severity]
		bg := severityBgColor[n.Severity]
		vid.DrawRoundedRect(
			25*menu.ratio,
			(stack+offset-46)*menu.ratio,
			lw+40*menu.ratio,
			70*menu.ratio,
			0.25,
			video.Color{R: float32(bg.R), G: float32(bg.G), B: float32(bg.B), A: fading},
		)
		vid.Font.SetColor(float32(fg.R), float32(fg.G), float32(fg.B), fading)
		vid.Font.Printf(
			45*menu.ratio,
			(stack+offset)*menu.ratio,
			0.5*menu.ratio,
			n.Message,
		)
		stack += h + offset
	}
}
