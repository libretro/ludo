package video

// Color is an RGBA type that we use in the menu
type Color struct {
	R, G, B, A float32
}

var white = Color{R: 1, G: 1, B: 1, A: 1}
var black = Color{R: 0, G: 0, B: 0, A: 1}
var deepGrey = Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
var grey = Color{R: 0.5, G: 0.5, B: 0.5, A: 1}
var lightGrey = Color{R: 0.75, G: 0.75, B: 0.75, A: 1}
var lightBlue = Color{R: 0.8784, G: 1, B: 1, A: 1}
var deepBlue = Color{R: 0.1, G: 0.1, B: 0.4, A: 0.85}

// Theme interface used by video Struct to defined theme
type Theme interface {
	GetMain() Color
	GetMainInverted() Color
	GetHintbarSecondary() Color
	GetGrey() Color
	GetHintbarPrimary() Color
	GetHightlightCursor() Color
	GetHCLHexagon() float64
	GetHexagonIconColor() Color
}

// LightTheme defined the lightTheme struct
type LightTheme struct{}

// DarkTheme defined the lightTheme struct
type DarkTheme struct{}

// GetMain return color of the theme
func (theme *LightTheme) GetMain() Color {
	return white
}

// GetMainInverted return color of the theme
func (theme *LightTheme) GetMainInverted() Color {
	return black
}

// GetHintbarSecondary return color of the theme
func (theme *LightTheme) GetHintbarSecondary() Color {
	return deepGrey
}

// GetGrey return color of the theme
func (theme *LightTheme) GetGrey() Color {
	return grey
}

// GetHintbarPrimary return hintBar color theme
func (theme *LightTheme) GetHintbarPrimary() Color {
	return lightGrey
}

// GetHightlightCursor return color of the cursor
func (theme *LightTheme) GetHightlightCursor() Color {
	return lightBlue
}

// GetHCLHexagon return color of the theme
func (theme *LightTheme) GetHCLHexagon() float64 {
	return 0
}

// GetHexagonIconColor return color of the theme
func (theme *LightTheme) GetHexagonIconColor() Color {
	return white
}

// GetMain return color of the theme
func (theme *DarkTheme) GetMain() Color {
	return black
}

// GetMainInverted return color of the theme
func (theme *DarkTheme) GetMainInverted() Color {
	return white
}

// GetHintbarSecondary return color of the theme
func (theme *DarkTheme) GetHintbarSecondary() Color {
	return lightGrey
}

// GetGrey return color of the theme
func (theme *DarkTheme) GetGrey() Color {
	return grey
}

// GetHintbarPrimary return hintBar color theme
func (theme *DarkTheme) GetHintbarPrimary() Color {
	return deepGrey
}

// GetHightlightCursor return color of the cursor
func (theme *DarkTheme) GetHightlightCursor() Color {
	return deepBlue
}

// GetHCLHexagon return color of the theme
func (theme *DarkTheme) GetHCLHexagon() float64 {
	return 180
}

// GetHexagonIconColor return color of the theme
func (theme *DarkTheme) GetHexagonIconColor() Color {
	return white
}
