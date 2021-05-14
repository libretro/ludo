// Package state holds the global state of the app. It is a separate package
// so we can import it from anywhere.
package state

import (
	"github.com/libretro/ludo/dat"
	"github.com/libretro/ludo/libretro"
)

var (
	Core        *libretro.Core // Core is the current libretro core, if any is loaded
	CoreRunning bool           // CoreRunning is true if a game or a gameless core is loaded
	MenuActive  bool           // MenuActive is whether to display the menu layer
	Verbose     bool           // Verbose will output more logs
	CorePath    string         // CorePath is the path of the current libretro core
	GamePath    string         // GamePath is the path of the current game
	DB          dat.DB         // DB is the game database loaded on startup
	LudOS       bool           // LudOS is a flag to run Ludo as a unix desktop environment
	FastForward bool           // FastForward will run the core as fast as possible
)
