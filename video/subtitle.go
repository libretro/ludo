package video

import (
	"strings"
	"time"
)

// RenderSubtitle draws the current subtitle, if any.
func (video *Video) RenderSubtitle() {
	if video == nil || video.Font == nil || video.Window == nil {
		return
	}

	if video.subtitleText == "" || time.Now().After(video.subtitleUntil) {
		video.subtitleText = ""
		return
	}

	fbw, fbh := video.Window.GetFramebufferSize()
	if fbw == 0 || fbh == 0 {
		return
	}

	lines := strings.Split(video.subtitleText, "\n")
	ratio := float32(fbw) / 1920
	scale := 0.6 * ratio
	padding := 16 * ratio
	lineHeight := 64 * ratio
	baseY := float32(fbh) - 80*ratio

	video.Font.UpdateResolution(fbw, fbh)

	previousColor := video.Font.color
	defer video.Font.SetColor(previousColor)

	var maxWidth float32
	for _, line := range lines {
		w := video.Font.Width(scale, line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	topY := baseY - lineHeight*float32(len(lines)-1)
	bgX := float32(fbw)/2 - maxWidth/2 - padding
	bgW := maxWidth + padding*2
	bgH := lineHeight*float32(len(lines)) + padding*2
	bgY := topY - padding

	video.DrawRect(bgX, bgY, bgW, bgH, 12*ratio, Color{0, 0, 0, 0.7})
	video.Font.SetColor(Color{1, 1, 1, 1})

	for i, line := range lines {
		y := topY + lineHeight*float32(i)
		x := float32(fbw)/2 - video.Font.Width(scale, line)/2
		video.Font.Printf(x, y, scale, line)
	}
}
