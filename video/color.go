package video

// Color is an RGBA type that we use in the menu
type Color struct {
	R, G, B, A float32
}

// Theme interface used by video Struct to defined theme
type Theme interface {
	GetMain() Color
	GetMainInverted() Color
	GetDeepGrey() Color
	GetGrey() Color
	GetLightGrey() Color
	GetHightlightCursor() Color
	GetHCLHexagon() float64
}

// LightTheme defined the lightTheme struct
type LightTheme struct {
	main             Color
	mainInverted     Color
	deepGrey         Color
	grey             Color
	lightGrey        Color
	hightlightCursor Color
}

// DarkTheme defined the lightTheme struct
type DarkTheme struct {
	main             Color
	mainInverted     Color
	deepGrey         Color
	grey             Color
	lightGrey        Color
	hightlightCursor Color
}

// NewLightTheme create the lightTheme
func NewLightTheme() *LightTheme {
	theme := LightTheme{}
	theme.main = Color{R: 1, G: 1, B: 1, A: 1}
	theme.mainInverted = Color{R: 0, G: 0, B: 0, A: 1}
	theme.deepGrey = Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	theme.grey = Color{R: 0.5, G: 0.5, B: 0.5, A: 1}
	theme.lightGrey = Color{R: 0.75, G: 0.75, B: 0.75, A: 1}
	theme.hightlightCursor = Color{R: 0.8784, G: 1, B: 1, A: 1}
	return &theme
}

// NewDarkTheme create the darkTheme
func NewDarkTheme() *DarkTheme {
	theme := DarkTheme{}
	theme.main = Color{R: 0, G: 0, B: 0, A: 1}
	theme.mainInverted = Color{R: 1, G: 1, B: 1, A: 1}
	theme.deepGrey = Color{R: 0.75, G: 0.75, B: 0.75, A: 1}
	theme.grey = Color{R: 0.5, G: 0.5, B: 0.5, A: 1}
	theme.lightGrey = Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	theme.hightlightCursor = Color{R: 0.1, G: 0.1, B: 0.4, A: 0.85}
	return &theme
}

// GetMain return main color of the theme
func (theme *LightTheme) GetMain() Color {
	return theme.main
}

// GetMainInverted return main color of the theme
func (theme *LightTheme) GetMainInverted() Color {
	return theme.mainInverted
}

// GetDeepGrey return main color of the theme
func (theme *LightTheme) GetDeepGrey() Color {
	return theme.deepGrey
}

// GetGrey return main color of the theme
func (theme *LightTheme) GetGrey() Color {
	return theme.grey
}

// GetLightGrey return main color of the theme
func (theme *LightTheme) GetLightGrey() Color {
	return theme.lightGrey
}

// GetHightlightCursor return main color of the theme
func (theme *LightTheme) GetHightlightCursor() Color {
	return theme.hightlightCursor
}

// GetHCLHexagon return main color of the theme
func (theme *LightTheme) GetHCLHexagon() float64 {
	return 0
}

// GetMain return main color of the theme
func (theme *DarkTheme) GetMain() Color {
	return theme.main
}

// GetMainInverted return main color of the theme
func (theme *DarkTheme) GetMainInverted() Color {
	return theme.mainInverted
}

// GetDeepGrey return main color of the theme
func (theme *DarkTheme) GetDeepGrey() Color {
	return theme.deepGrey
}

// GetGrey return main color of the theme
func (theme *DarkTheme) GetGrey() Color {
	return theme.grey
}

// GetLightGrey return main color of the theme
func (theme *DarkTheme) GetLightGrey() Color {
	return theme.lightGrey
}

// GetHightlightCursor return main color of the theme
func (theme *DarkTheme) GetHightlightCursor() Color {
	return theme.hightlightCursor
}

// GetHCLHexagon return main color of the theme
func (theme *DarkTheme) GetHCLHexagon() float64 {
	return 180
}
