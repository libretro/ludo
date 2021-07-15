// Package state holds the global state of the app. It is a separate package
// so we can import it from anywhere.
package state

import (
	"github.com/libretro/ludo/dat"
	"github.com/libretro/ludo/libretro"
)

// Core is the current libretro core, if any is loaded
var Core *libretro.Core

// CoreRunning is true if a game or a gameless core is loaded
var CoreRunning bool

// MenuActive is whether to display the menu layer
var MenuActive bool

// Verbose will output more logs
var Verbose bool

// CorePath is the path of the current libretro core
var CorePath string

// GamePath is the path of the current game
var GamePath string

// DB is the game database loaded on startup
var DB dat.DB

// LudOS is whether run Ludo as a unix desktop environment
var LudOS bool

// FastForward will run the core as fast as possible
var FastForward bool

// Tick is the current frame, used for rollback networking
var Tick int64

// ForcePause is used to debug rollback networking
var ForcePause bool
