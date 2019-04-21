package menu

import (
	"github.com/libretro/ludo/notifications"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/video"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var severityFgColor = map[ntf.Severity]colorful.Color{
	ntf.Error:   colorful.Hcl(40, 0.75, 0.85),
	ntf.Warning: colorful.Hcl(90, 0.75, 0.85),
	ntf.Success: colorful.Hcl(120, 0.75, 0.85),
	ntf.Info:    colorful.Hcl(230, 0.75, 0.85),
}

var severityBgColor = map[ntf.Severity]colorful.Color{
	ntf.Error:   colorful.Hcl(40, 0.65, 0.1),
	ntf.Warning: colorful.Hcl(90, 0.65, 0.1),
	ntf.Success: colorful.Hcl(120, 0.65, 0.1),
	ntf.Info:    colorful.Hcl(230, 0.65, 0.1),
}

// RenderNotifications draws the list of notification messages on the viewport
func (m *Menu) RenderNotifications() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)
	var h float32 = 75
	stack := h
	for _, n := range notifications.List() {
		fading := float32(n.Duration) * 4
		if fading > 1 {
			fading = 1
		}
		offset := fading*h - h
		lw := vid.Font.Width(0.5*m.ratio, n.Message)
		fg := severityFgColor[n.Severity]
		bg := severityBgColor[n.Severity]
		vid.DrawRoundedRect(
			25*m.ratio,
			(stack+offset-46)*m.ratio,
			lw+40*m.ratio,
			70*m.ratio,
			0.25,
			video.Color{R: float32(bg.R), G: float32(bg.G), B: float32(bg.B), A: fading},
		)
		vid.Font.SetColor(float32(fg.R), float32(fg.G), float32(fg.B), fading)
		vid.Font.Printf(
			45*m.ratio,
			(stack+offset)*m.ratio,
			0.5*m.ratio,
			n.Message,
		)
		stack += h + offset
	}
}
