package menu

import (
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
	for _, n := range ntf.List() {
		offset := minf32(n.Duration*4, 1)*h - h
		lw := vid.Font.Width(0.4*menu.ratio, n.Message)

		Toast(&Props{
			X:            25 * menu.ratio,
			Y:            (stack + offset - 46) * menu.ratio,
			Width:        lw + 70*menu.ratio + 20*menu.ratio,
			Height:       70 * menu.ratio,
			BorderRadius: 0.25,
		}, n)()

		stack += h + offset
	}
}

// Toast can render a notification
func Toast(props *Props, n *ntf.Notification) func() {
	fg := severityFgColor[n.Severity]
	bg := severityBgColor[n.Severity]
	alpha := minf32(n.Duration*4, 1)
	props.Color = video.Color{R: float32(bg.R), G: float32(bg.G), B: float32(bg.B), A: alpha}

	return HBox(props,
		Image(&Props{
			Width:  props.Height,
			Height: props.Height,
			Scale:  1,
			Color:  video.Color{R: float32(fg.R), G: float32(fg.G), B: float32(fg.B), A: alpha},
		}, menu.icons["core-infos"]),
		Label(&Props{
			Height: props.Height,
			Scale:  0.4 * menu.ratio,
			Color:  video.Color{R: float32(fg.R), G: float32(fg.G), B: float32(fg.B), A: alpha},
		}, n.Message),
	)
}
