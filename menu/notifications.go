package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/video"
)

var severityFgColor = map[ntf.Severity]video.Color{
	ntf.Error:   lightDanger,
	ntf.Warning: lightWarning,
	ntf.Success: lightSuccess,
	ntf.Info:    lightInfo,
}

var severityBgColor = map[ntf.Severity]video.Color{
	ntf.Error:   darkDanger,
	ntf.Warning: darkWarning,
	ntf.Success: darkSuccess,
	ntf.Info:    darkInfo,
}

// RenderNotifications draws the list of notification messages on the viewport
func (m *Menu) RenderNotifications() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)
	var h float32 = 75
	stack := h
	for _, n := range ntf.List() {
		fading := n.Duration * 4
		if fading > 1 {
			fading = 1
		}
		offset := fading*h - h
		lw := vid.Font.Width(0.5*m.ratio, n.Message)
		fg := severityFgColor[n.Severity]
		bg := severityBgColor[n.Severity]
		vid.DrawRect(
			25*m.ratio,
			(stack+offset-46)*m.ratio,
			lw+40*m.ratio,
			70*m.ratio,
			0.25,
			bg.Alpha(fading),
		)
		vid.Font.SetColor(fg.Alpha(fading))
		vid.Font.Printf(
			45*m.ratio,
			(stack+offset)*m.ratio,
			0.5*m.ratio,
			n.Message,
		)
		stack += h + offset
	}
}
