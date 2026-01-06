package video

import (
	"strings"
	"time"
)

// SetSubtitle displays text at the bottom of the screen for a given duration.
func (video *Video) SetSubtitle(text string, duration time.Duration) {
	video.subtitleText = strings.TrimSpace(text)
	if video.subtitleText == "" {
		video.subtitleUntil = time.Time{}
		return
	}

	video.subtitleUntil = time.Now().Add(duration)
}

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
	bgW := maxWidth + padding*2
	minBgW := float32(fbw) * 0.75
	if bgW < minBgW {
		bgW = minBgW
	}
	bgX := (float32(fbw) - bgW) / 2
	bgH := lineHeight*float32(len(lines)) + padding*2
	bgY := topY - padding

	bgColor := Color{0, 0, 0, 0.7}
	borderColor := Color{1, 1, 1, 0.08}

	video.DrawRect(bgX, bgY, bgW, bgH, 12*ratio, bgColor)
	video.DrawBorder(bgX, bgY, bgW, bgH, 2*ratio, borderColor)
	video.Font.SetColor(Color{1, 1, 1, 1})

	for i, line := range lines {
		y := topY + lineHeight*float32(i)
		x := float32(fbw)/2 - video.Font.Width(scale, line)/2
		video.Font.Printf(x, y, scale, line)
	}
}
