// Package state holds the global state of the app. It is a separate package
// so we can import it from anywhere.
package state

import (
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/rdb"
	"github.com/libretro/ludo/tasks"
)

// State is a type for the global state of the app
type State struct {
	Core        libretro.Core // Current libretro core
	FrameTimeCb libretro.FrameTimeCallback
	AudioCb     libretro.AudioCallback
	CoreRunning bool
	MenuActive  bool // When set to true, will display the menu layer
	Verbose     bool
	CorePath    string // Path of the current libretro core
	GamePath    string // Path of the current game
	DB          rdb.DB
	Tasks       []tasks.Task
}

// Global state
var Global State
