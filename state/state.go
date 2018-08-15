package state

import "github.com/libretro/go-playthemall/libretro"

// global state
var Global struct {
	Core        libretro.Core
	FrameTimeCb libretro.FrameTimeCallback
	AudioCb     libretro.AudioCallback
	CoreRunning bool
	MenuActive  bool
	Verbose     bool
	CorePath    string
	GamePath    string
}
