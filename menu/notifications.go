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
	stack := float32(0)
	for _, n := range ntf.List() {
		offset := minf32(n.Duration*4, 1) * 80
		lw := vid.Font.Width(0.4*menu.ratio, n.Message)

		Toast(&Props{
			X:            25 * menu.ratio,
			Y:            (stack + offset - 80 + 25) * menu.ratio,
			Width:        lw + 70*menu.ratio + 20*menu.ratio,
			Height:       70 * menu.ratio,
			BorderRadius: 0.25,
		}, n)()

		stack += offset
	}
}

// Toast can render a notification
func Toast(props *Props, n *ntf.Notification) func() {
	alpha := minf32(n.Duration*4, 1)
	fg := severityFgColor[n.Severity].Alpha(alpha)
	bg := severityBgColor[n.Severity].Alpha(alpha)
	props.Color = bg

	return HBox(props,
		Image(&Props{
			Width:  props.Height,
			Height: props.Height,
			Scale:  1,
			Color:  fg,
		}, menu.icons["core-infos"]),
		Label(&Props{
			Height: props.Height,
			Scale:  0.4 * menu.ratio,
			Color:  fg,
		}, n.Message),
	)
}
