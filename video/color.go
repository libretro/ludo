package video

import "github.com/libretro/ludo/settings"

// Color is an RGBA type that we use in the menu
type Color struct {
	R, G, B, A float32
}

var lightTheme = map[string]Color{
	"main":              {R: 1, G: 1, B: 1, A: 0.85},
	"main-inverted":     {R: 0, G: 0, B: 0, A: 0.85},
	"deep-grey":         {R: 0.25, G: 0.25, B: 0.25, A: 1},
	"grey":              {R: 0.5, G: 0.5, B: 0.5, A: 1},
	"light-grey":        {R: 0.75, G: 0.75, B: 0.75, A: 1},
	"hightlight-cursor": {R: 0.8784, G: 1, B: 1, A: 0.85},
}

var darkTheme = map[string]Color{
	"main":              {R: 0, G: 0, B: 0, A: 0.85},
	"main-inverted":     {R: 1, G: 1, B: 1, A: 0.85},
	"deep-grey":         {R: 0.25, G: 0.25, B: 0.25, A: 1},
	"grey":              {R: 0.5, G: 0.5, B: 0.5, A: 1},
	"light-grey":        {R: 0.75, G: 0.75, B: 0.75, A: 1},
	"hightlight-cursor": {R: 0.1, G: 0.1, B: 0.4, A: 0.85},
}

// GetThemeColor colors by theme using the static map define in color.go
func GetThemeColor(colorName string, alpha float32) Color {
	var c = Color{}
	if settings.Current.VideoDarkMode {
		c = darkTheme[colorName]
	} else {
		c = lightTheme[colorName]
	}
	c.A = alpha
	return c
}
