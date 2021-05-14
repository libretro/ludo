// Package state holds the global state of the app. It is a separate package
// so we can import it from anywhere.
package state

import (
	"github.com/libretro/ludo/dat"
	"github.com/libretro/ludo/libretro"
)

var (
	Core        *libretro.Core // Current libretro core
	CoreRunning bool           // Should be true if a game or a gameless core is loaded
	MenuActive  bool           // When set to true, will display the menu layer
	Verbose     bool           // When set to true, will output a lots of logs
	CorePath    string         // Path of the current libretro core
	GamePath    string         // Path of the current game
	DB          dat.DB         // The game database loaded on startup
	LudOS       bool           // Run Ludo as a unix desktop environment
	FastForward bool           // Run the core as fast as possible
)
