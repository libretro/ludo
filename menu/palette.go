package menu

import (
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

var white = video.Color{R: 1, G: 1, B: 1, A: 1}
var black = video.Color{R: 0, G: 0, B: 0, A: 1}

var orange = video.Color{R: 0.8, G: 0.4, B: 0.1, A: 1}
var cyan = video.Color{R: 0.8784, G: 1, B: 1, A: 1}
var darkBlue = video.Color{R: 0.1, G: 0.1, B: 0.4, A: 1}

var lightGrey = video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1}
var mediumGrey = video.Color{R: 0.5, G: 0.5, B: 0.5, A: 1}
var darkGrey = video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
var darkerGrey = video.Color{R: 0.10, G: 0.10, B: 0.10, A: 1}

var darkInfo = video.Color{R: 0.04, G: 0.36, B: 0.46, A: 1}
var lightInfo = video.Color{R: 0.53, G: 0.89, B: 1.00, A: 1}

var darkSuccess = video.Color{R: 0.15, G: 0.46, B: 0.04, A: 1}
var lightSuccess = video.Color{R: 0.65, G: 1.00, B: 0.53, A: 1}

var darkDanger = video.Color{R: 0.46, G: 0.04, B: 0.04, A: 1}
var lightDanger = video.Color{R: 1.00, G: 0.53, B: 0.53, A: 1}

var darkWarning = video.Color{R: 0.47, G: 0.40, B: 0.04, A: 1}
var lightWarning = video.Color{R: 1.00, G: 0.92, B: 0.53, A: 1}

var bgColor = white
var cursorBg = cyan
var textColor = black

// UpdatePalette updates the color palette to honor the dark theme
func (m *Menu) UpdatePalette() {
	bgColor = white
	cursorBg = cyan
	textColor = black

	if state.Global.CoreRunning || settings.Current.VideoDarkMode {
		bgColor = darkerGrey
		cursorBg = darkBlue
		textColor = white
	}
}
