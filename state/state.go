// Package state holds the global state of the app. It is a separate package
// so we can import it from anywhere.
package state

import (
	"sync"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/rdb"
)

// State is a type for the global state of the app
type State struct {
	Core              *libretro.Core // Current libretro core
	CoreRunning       bool           // Should be true if a game or a gameless core is loaded
	MenuActive        bool           // When set to true, will display the menu layer
	Verbose           bool           // When set to true, will output a lots of logs
	CorePath          string         // Path of the current libretro core
	GamePath          string         // Path of the current game
	DB                rdb.DB         // The game database loaded on startup
	LudOS             bool           // Run Ludo as a unix desktop environment
	FastForward       bool           // Run the core as fast as possible
	Tick              int64
	LastConfirmedTick int64
	SyncedLastUpdate  bool
	ForcePause        bool

	sync.Mutex
}

// Global state
var Global State
