package menu

import (
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

var white = video.Color{R: 250.0/255.0, G: 244.0/255.0, B: 237.0/255.0, A: 1}
var black = video.Color{R: 87.0/255.0, G: 82.0/255.0, B: 121.0/255.0, A: 1}

var foam = video.Color{R: 86.0/255.0, G: 148.0/255.0, B: 159.0/255.0, A: 1}
var overlay = video.Color{R: 242.0/255.0, G: 233.0/255.0, B: 222.0/255.0, A: 1}

var highlightMed = video.Color{R: 223.0/255.0, G: 218.0/255.0, B: 217.0/255.0, A: 1}
var muted = video.Color{R: 152.0/255.0, G: 147.0/255.0, B: 165.0/255.0, A: 1}

var darkInfo = video.Color{R: 144.0/255.0, G: 122.0/255.0, B: 169.0/255.0, A: 1}
var lightInfo = white

var darkSuccess = foam
var lightSuccess = white

var darkDanger = video.Color{R: 180.0/255.0, G: 99.0/255.0, B: 122.0/255.0, A: 1}
var lightDanger = white

var darkWarning = video.Color{R: 234.0/255.0, G: 157.0/255.0, B: 52.0/255.0, A: 1}
var lightWarning = white

var bgColor = white
var cursorBg = overlay
var textColor = black
var mutedTextColor = muted
var textShadowColor = overlay
var sepColor = highlightMed
var hintTextColor = black
var hintBgColor = white
var titleTextColor = foam
var dialogTextColor = black

// UpdatePalette updates the color palette to honor the dark theme
func (m *Menu) UpdatePalette() {
	bgColor = white
	cursorBg = overlay
	textColor = black
	textShadowColor = overlay
	sepColor = highlightMed
	hintTextColor = black
	hintBgColor = white
	dialogTextColor = black

	if state.CoreRunning || settings.Current.VideoDarkMode {
		bgColor = black
		cursorBg = black
		textColor = white
		textShadowColor = black
		sepColor = white
		hintTextColor = white
		hintBgColor = black
		dialogTextColor = black
	}
}
