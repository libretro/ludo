package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

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

var infoBgColor video.Color
var infoTextColor video.Color
var successBgColor video.Color
var successTextColor video.Color
var dangerBgColor video.Color
var dangerTextColor video.Color
var warningBgColor video.Color
var warningTextColor video.Color

var tabColors []video.Color
var tabIconColor video.Color

func applyDraculaTheme() {
	var dracula = map[string]video.Color{
		"background": video.ColorFromHex("#282A36"),
		"foreground": video.ColorFromHex("#F8F8F2"),
		"selection":  video.ColorFromHex("#44475A"),
		"comment":    video.ColorFromHex("#6272A4"),
		"red":        video.ColorFromHex("#FF5555"),
		"orange":     video.ColorFromHex("#FFB86C"),
		"yellow":     video.ColorFromHex("#F1FA8C"),
		"green":      video.ColorFromHex("#50FA7B"),
		"purple":     video.ColorFromHex("#BD93F9"),
		"cyan":       video.ColorFromHex("#8BE9FD"),
		"pink":       video.ColorFromHex("#FF79C6"),
	}

	white = dracula["foreground"]
	black = dracula["background"]
	bgColor = dracula["background"]
	cursorBgColor = dracula["selection"]
	textColor = dracula["foreground"]
	mutedTextColor = dracula["comment"]
	textShadowColor = dracula["selection"]
	sepColor = dracula["selection"]
	hintTextColor = dracula["foreground"]
	hintBgColor = dracula["background"]
	titleTextColor = dracula["pink"]
	dialogTextColor = dracula["foreground"]

	infoBgColor = dracula["selection"]
	infoTextColor = dracula["purple"]
	successBgColor = dracula["selection"]
	successTextColor = dracula["green"]
	dangerBgColor = dracula["selection"]
	dangerTextColor = dracula["red"]
	warningBgColor = dracula["selection"]
	warningTextColor = dracula["orange"]

	tabColors = []video.Color{
		dracula["purple"],
		dracula["red"],
		dracula["pink"],
		dracula["orange"],
		dracula["green"],
		dracula["cyan"],
	}

	tabIconColor = dracula["foreground"]
}

// rosePine color palette
// This is the dark theme color palette used in the menu
func applyRosePineTheme() {
	var rosePine = map[string]video.Color{
		"base":         video.ColorFromHex("#191724"),
		"text":         video.ColorFromHex("#e0def4"),
		"overlay":      video.ColorFromHex("#26233a"),
		"muted":        video.ColorFromHex("#6e6a86"),
		"love":         video.ColorFromHex("#eb6f92"),
		"gold":         video.ColorFromHex("#f6c177"),
		"rose":         video.ColorFromHex("#ebbcba"),
		"pine":         video.ColorFromHex("#31748f"),
		"foam":         video.ColorFromHex("#9ccfd8"),
		"iris":         video.ColorFromHex("#c4a7e7"),
		"highlightMed": video.ColorFromHex("#403d52"),
	}

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
	titleTextColor = rosePine["pine"]
	dialogTextColor = rosePine["base"]

	infoBgColor = rosePine["iris"]
	infoTextColor = rosePine["base"]
	successBgColor = rosePine["foam"]
	successTextColor = rosePine["base"]
	dangerBgColor = rosePine["love"]
	dangerTextColor = rosePine["base"]
	warningBgColor = rosePine["gold"]
	warningTextColor = rosePine["base"]

	tabColors = []video.Color{
		rosePine["iris"],
		rosePine["love"],
		rosePine["rose"],
		rosePine["gold"],
		rosePine["foam"],
		rosePine["pine"],
	}

	tabIconColor = rosePine["text"]
}

// rosePineDawn color palette
// This is the light theme color palette used in the menu
func applyRosePineDawnTheme() {
	var rosePineDawn = map[string]video.Color{
		"base":         video.ColorFromHex("#faf4ed"),
		"text":         video.ColorFromHex("#575279"),
		"overlay":      video.ColorFromHex("#f2e9e1"),
		"muted":        video.ColorFromHex("#9893a5"),
		"love":         video.ColorFromHex("#b4637a"),
		"gold":         video.ColorFromHex("#ea9d34"),
		"rose":         video.ColorFromHex("#d7827e"),
		"pine":         video.ColorFromHex("#286983"),
		"foam":         video.ColorFromHex("#56949f"),
		"iris":         video.ColorFromHex("#907aa9"),
		"highlightMed": video.ColorFromHex("#dfdad9"),
	}

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

	infoBgColor = rosePineDawn["iris"]
	infoTextColor = rosePineDawn["base"]
	successBgColor = rosePineDawn["foam"]
	successTextColor = rosePineDawn["base"]
	dangerBgColor = rosePineDawn["love"]
	dangerTextColor = rosePineDawn["base"]
	warningBgColor = rosePineDawn["gold"]
	warningTextColor = rosePineDawn["base"]

	tabColors = []video.Color{
		rosePineDawn["iris"],
		rosePineDawn["love"],
		rosePineDawn["rose"],
		rosePineDawn["gold"],
		rosePineDawn["foam"],
		rosePineDawn["pine"],
	}

	tabIconColor = rosePineDawn["base"]
}

// UpdatePalette updates the color palette to honor the dark theme
func (m *Menu) UpdatePalette() {
	if state.CoreRunning || settings.Current.VideoDarkMode {
		applyRosePineTheme()
	} else {
		applyRosePineDawnTheme()
	}

	severityFgColor = map[ntf.Severity]video.Color{
		ntf.Error:   dangerTextColor,
		ntf.Warning: warningTextColor,
		ntf.Success: successTextColor,
		ntf.Info:    infoTextColor,
	}

	severityBgColor = map[ntf.Severity]video.Color{
			ntf.Error:   dangerBgColor,
			ntf.Warning: warningBgColor,
			ntf.Success: successBgColor,
			ntf.Info:    infoBgColor,
	}
}
