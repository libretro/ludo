package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

// rosePineDawn color palette
// This is the light theme color palette used in the menu
var rosePineDawn = map[string]video.Color{
	"base": video.Color{R: 250.0/255.0, G: 244.0/255.0, B: 237.0/255.0, A: 1},
	"text": video.Color{R: 87.0/255.0, G: 82.0/255.0, B: 121.0/255.0, A: 1},
	"overlay": video.Color{R: 242.0/255.0, G: 233.0/255.0, B: 222.0/255.0, A: 1},
	"muted": video.Color{R: 152.0/255.0, G: 147.0/255.0, B: 165.0/255.0, A: 1},
	"love": video.Color{R: 180.0/255.0, G: 99.0/255.0, B: 122.0/255.0, A: 1},
	"gold": video.Color{R: 234.0/255.0, G: 157.0/255.0, B: 52.0/255.0, A: 1},
	"rose": video.Color{R: 215.0/255.0, G: 130.0/255.0, B: 126.0/255.0, A: 1},
	"pine": video.Color{R: 40.0/255.0, G: 105.0/255.0, B: 131.0/255.0, A: 1},
	"foam": video.Color{R: 86.0/255.0, G: 148.0/255.0, B: 159.0/255.0, A: 1},
	"iris": video.Color{R: 144.0/255.0, G: 122.0/255.0, B: 169.0/255.0, A: 1},
	"highlightMed": video.Color{R: 223.0/255.0, G: 218.0/255.0, B: 217.0/255.0, A: 1},
}

// rosePine color palette
// This is the dark theme color palette used in the menu
var rosePine = map[string]video.Color{
	"base": video.Color{R: 25.0/255.0, G: 23.0/255.0, B: 36.0/255.0, A: 1},
	"text": video.Color{R: 224.0/255.0, G: 222.0/255.0, B: 244.0/255.0, A: 1},
	"overlay": video.Color{R: 38.0/255.0, G: 35.0/255.0, B: 58.0/255.0, A: 1},
	"muted": video.Color{R: 110.0/255.0, G: 106.0/255.0, B: 134.0/255.0, A: 1},
	"love": video.Color{R: 235.0/255.0, G: 111.0/255.0, B: 146.0/255.0, A: 1},
	"gold": video.Color{R: 246.0/255.0, G: 193.0/255.0, B: 119.0/255.0, A: 1},
	"rose": video.Color{R: 235.0/255.0, G: 188.0/255.0, B: 186.0/255.0, A: 1},
	"pine": video.Color{R: 49.0/255.0, G: 116.0/255.0, B: 143.0/255.0, A: 1},
	"foam": video.Color{R: 156.0/255.0, G: 207.0/255.0, B: 216.0/255.0, A: 1},
	"iris": video.Color{R: 196.0/255.0, G: 167.0/255.0, B: 231.0/255.0, A: 1},
	"highlightMed": video.Color{R: 64.0/255.0, G: 61.0/255.0, B: 82.0/255.0, A: 1},
}

// Global variables for the menu theme

var white video.Color
var black video.Color
var bgColor video.Color
var cursorBgColor video.Color
var textColor video.Color
var mutedTextColor video.Color
var textShadowColor video.Color
var sepColor video.Color
var hintTextColor video.Color
var hintBgColor video.Color
var titleTextColor video.Color
var dialogTextColor video.Color

var darkInfo video.Color
var lightInfo video.Color
var darkSuccess video.Color
var lightSuccess video.Color
var darkDanger video.Color
var lightDanger video.Color
var darkWarning video.Color
var lightWarning video.Color

var tabColors []video.Color
var tabIconColor video.Color

// UpdatePalette updates the color palette to honor the dark theme
func (m *Menu) UpdatePalette() {
	if state.CoreRunning || settings.Current.VideoDarkMode {
		white = rosePine["text"]
		black = rosePine["base"]
		bgColor = rosePine["base"]
		cursorBgColor = rosePine["overlay"]
		textColor = rosePine["text"]
		mutedTextColor = rosePine["muted"]
		textShadowColor = rosePine["overlay"]
		sepColor = rosePine["highlightMed"]
		hintTextColor = rosePine["text"]
		hintBgColor = rosePine["base"]
		titleTextColor = rosePine["foam"]
		dialogTextColor = rosePine["text"]

		darkInfo = rosePine["iris"]
		lightInfo = rosePine["base"]
		darkSuccess = rosePine["foam"]
		lightSuccess = rosePine["base"]
		darkDanger = rosePine["love"]
		lightDanger = rosePine["base"]
		darkWarning = rosePine["gold"]
		lightWarning = rosePine["base"]

		tabColors = []video.Color{
			rosePine["iris"],
			rosePine["love"],
			rosePine["rose"],
			rosePine["gold"],
			rosePine["foam"],
			rosePine["pine"],
		}

		tabIconColor = rosePine["text"]
	} else {
		white = rosePineDawn["base"]
		black = rosePineDawn["text"]
		bgColor = rosePineDawn["base"]
		cursorBgColor = rosePineDawn["overlay"]
		textColor = rosePineDawn["text"]
		mutedTextColor = rosePineDawn["muted"]
		textShadowColor = rosePineDawn["overlay"]
		sepColor = rosePineDawn["highlightMed"]
		hintTextColor = rosePineDawn["text"]
		hintBgColor = rosePineDawn["base"]
		titleTextColor = rosePineDawn["foam"]
		dialogTextColor = rosePineDawn["text"]

		darkInfo = rosePineDawn["iris"]
		lightInfo = rosePineDawn["base"]
		darkSuccess = rosePineDawn["foam"]
		lightSuccess = rosePineDawn["base"]
		darkDanger = rosePineDawn["love"]
		lightDanger = rosePineDawn["base"]
		darkWarning = rosePineDawn["gold"]
		lightWarning = rosePineDawn["base"]

		tabColors = []video.Color{
			rosePineDawn["iris"],
			rosePineDawn["love"],
			rosePineDawn["rose"],
			rosePineDawn["gold"],
			rosePineDawn["foam"],
			rosePineDawn["pine"],
		}

		tabIconColor = rosePineDawn["text"]
	}

	severityFgColor = map[ntf.Severity]video.Color{
		ntf.Error:   lightDanger,
		ntf.Warning: lightWarning,
		ntf.Success: lightSuccess,
		ntf.Info:    lightInfo,
	}

	severityBgColor = map[ntf.Severity]video.Color{
			ntf.Error:   darkDanger,
			ntf.Warning: darkWarning,
			ntf.Success: darkSuccess,
			ntf.Info:    darkInfo,
	}
}
