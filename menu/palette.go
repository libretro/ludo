package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"

	colorful "github.com/lucasb-eyer/go-colorful"
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

var tabHexaColors func(int) video.Color
var tabIconColors func(int) video.Color

func applyDefaultLightTheme() {
	var pal = map[string]video.Color{
		"background":  video.ColorFromHex("#FFFFFF"),
		"foreground":  video.ColorFromHex("#000000"),
		"selection":   {R: 0.8784, G: 1, B: 1, A: 1},
		"lightgrey":   video.ColorFromHex("#cccccc"),
		"mediumgrey":  video.ColorFromHex("#888888"),
		"darkgrey":    video.ColorFromHex("#555555"),
		"darkergrey":  video.ColorFromHex("#222222"),
		"infobg":      {R: 0.04, G: 0.36, B: 0.46, A: 1},
		"infotext":    {R: 0.53, G: 0.89, B: 1.00, A: 1},
		"successbg":   {R: 0.15, G: 0.46, B: 0.04, A: 1},
		"successtext": {R: 0.65, G: 1.00, B: 0.53, A: 1},
		"dangerbg":    {R: 0.46, G: 0.04, B: 0.04, A: 1},
		"dangertext":  {R: 1.00, G: 0.53, B: 0.53, A: 1},
		"warningbg":   {R: 0.47, G: 0.40, B: 0.04, A: 1},
		"warningtext": {R: 1.00, G: 0.92, B: 0.53, A: 1},
		"orange":      {R: 0.760506, G: 0.4733901, B: 0.3095738, A: 1},
	}

	white = pal["background"]
	black = pal["foreground"]
	bgColor = pal["background"]
	cursorBgColor = pal["selection"]
	textColor = pal["foreground"]
	mutedTextColor = pal["mediumgrey"]
	textShadowColor = pal["lightgrey"]
	sepColor = pal["lightgrey"]
	hintTextColor = pal["darkergrey"]
	hintBgColor = pal["background"]
	titleTextColor = pal["orange"]
	dialogTextColor = pal["darkergrey"]

	infoBgColor = pal["infobg"]
	infoTextColor = pal["infotext"]
	successBgColor = pal["successbg"]
	successTextColor = pal["successtext"]
	dangerBgColor = pal["dangerbg"]
	dangerTextColor = pal["dangertext"]
	warningBgColor = pal["warningbg"]
	warningTextColor = pal["warningtext"]

	tabHexaColors = func(i int) video.Color {
		cf := colorful.Hcl(float64(i)*20, 0.5, 0.5)
		return video.Color{R: float32(cf.R), G: float32(cf.B), B: float32(cf.G), A: 1}
	}

	tabIconColors = func(i int) video.Color {
		return pal["background"]
	}
}

func applyDefaultDarkTheme() {
	var pal = map[string]video.Color{
		"background":  video.ColorFromHex("#000000"),
		"foreground":  video.ColorFromHex("#FFFFFF"),
		"selection":   {R: 0.8784, G: 1, B: 1, A: 1},
		"lightgrey":   video.ColorFromHex("#cccccc"),
		"mediumgrey":  video.ColorFromHex("#888888"),
		"darkgrey":    video.ColorFromHex("#555555"),
		"darkergrey":  video.ColorFromHex("#222222"),
		"infobg":      {R: 0.04, G: 0.36, B: 0.46, A: 1},
		"infotext":    {R: 0.53, G: 0.89, B: 1.00, A: 1},
		"successbg":   {R: 0.15, G: 0.46, B: 0.04, A: 1},
		"successtext": {R: 0.65, G: 1.00, B: 0.53, A: 1},
		"dangerbg":    {R: 0.46, G: 0.04, B: 0.04, A: 1},
		"dangertext":  {R: 1.00, G: 0.53, B: 0.53, A: 1},
		"warningbg":   {R: 0.47, G: 0.40, B: 0.04, A: 1},
		"warningtext": {R: 1.00, G: 0.92, B: 0.53, A: 1},
		"orange":      {R: 0.760506, G: 0.4733901, B: 0.3095738, A: 1},
	}

	white = pal["foreground"]
	black = pal["background"]
	bgColor = pal["background"]
	cursorBgColor = pal["darkergrey"]
	textColor = pal["foreground"]
	mutedTextColor = pal["mediumgrey"]
	textShadowColor = pal["darkergrey"]
	sepColor = pal["darkergrey"]
	hintTextColor = pal["lightgrey"]
	hintBgColor = pal["background"]
	titleTextColor = pal["orange"]
	dialogTextColor = pal["darkergrey"]

	infoBgColor = pal["infobg"]
	infoTextColor = pal["infotext"]
	successBgColor = pal["successbg"]
	successTextColor = pal["successtext"]
	dangerBgColor = pal["dangerbg"]
	dangerTextColor = pal["dangertext"]
	warningBgColor = pal["warningbg"]
	warningTextColor = pal["warningtext"]

	tabHexaColors = func(i int) video.Color {
		cf := colorful.Hcl(float64(i)*20, 0.5, 0.5)
		return video.Color{R: float32(cf.R), G: float32(cf.B), B: float32(cf.G), A: 1}
	}

	tabIconColors = func(_ int) video.Color {
		return pal["darkergrey"]
	}
}

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
	mutedTextColor = dracula["purple"]
	textShadowColor = dracula["selection"]
	sepColor = dracula["selection"]
	hintTextColor = dracula["foreground"]
	hintBgColor = dracula["background"]
	titleTextColor = dracula["purple"]
	dialogTextColor = dracula["background"]

	infoBgColor = dracula["selection"]
	infoTextColor = dracula["purple"]
	successBgColor = dracula["selection"]
	successTextColor = dracula["green"]
	dangerBgColor = dracula["selection"]
	dangerTextColor = dracula["red"]
	warningBgColor = dracula["selection"]
	warningTextColor = dracula["orange"]

	tabHexaColors = func(i int) video.Color {
		return []video.Color{
			dracula["purple"],
			dracula["red"],
			dracula["pink"],
			dracula["orange"],
			dracula["green"],
			dracula["cyan"],
		}[i%6]
	}

	tabIconColors = func(_ int) video.Color {
		return dracula["selection"]
	}
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

	tabHexaColors = func(i int) video.Color {
		return[]video.Color{
			rosePine["iris"],
			rosePine["love"],
			rosePine["rose"],
			rosePine["gold"],
			rosePine["foam"],
			rosePine["pine"],
		}[i%6]
	}

	tabIconColors = func(_ int) video.Color {
		return rosePine["overlay"]
	}
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

	tabHexaColors = func(i int) video.Color {
		return []video.Color{
			rosePineDawn["iris"],
			rosePineDawn["love"],
			rosePineDawn["rose"],
			rosePineDawn["gold"],
			rosePineDawn["foam"],
			rosePineDawn["pine"],
		}[i%6]
	}

	tabIconColors = func(_ int) video.Color {
		return rosePineDawn["base"]
	}
}

// UpdatePalette updates the color palette to honor the dark theme
func (m *Menu) UpdatePalette() {
	if state.CoreRunning || settings.Current.VideoDarkMode {
		switch settings.Current.VideoTheme {
			case "Rose Pine": applyRosePineTheme()
			case "Dracula": applyDraculaTheme()
			default: applyDefaultDarkTheme()
		}
	} else {
		switch settings.Current.VideoTheme {
			case "Rose Pine": applyRosePineDawnTheme()
			case "Dracula": applyDraculaTheme()
			default: applyDefaultLightTheme()
		}
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
