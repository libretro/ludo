package menu

import "github.com/libretro/ludo/video"

var white = video.Color{R: 1, G: 1, B: 1, A: 1}
var black = video.Color{R: 0, G: 0, B: 0, A: 1}
var blue = video.Color{R: 0.129, G: 0.441, B: 0.684, A: 1}
var lightGrey = video.Color{R: 0.85, G: 0.85, B: 0.85, A: 1}
var mediumGrey = video.Color{R: 0.56, G: 0.56, B: 0.56, A: 1}
var darkGrey = video.Color{R: 0.28, G: 0.28, B: 0.28, A: 1}

var darkInfo = video.Color{R: 0.04, G: 0.36, B: 0.46, A: 1}
var lightInfo = video.Color{R: 0.53, G: 0.89, B: 1.00, A: 1}

var darkSuccess = video.Color{R: 0.15, G: 0.46, B: 0.04, A: 1}
var lightSuccess = video.Color{R: 0.65, G: 1.00, B: 0.53, A: 1}

var darkDanger = video.Color{R: 0.46, G: 0.04, B: 0.04, A: 1}
var lightDanger = video.Color{R: 1.00, G: 0.53, B: 0.53, A: 1}

var darkWarning = video.Color{R: 0.47, G: 0.40, B: 0.04, A: 1}
var lightWarning = video.Color{R: 1.00, G: 0.92, B: 0.53, A: 1}
