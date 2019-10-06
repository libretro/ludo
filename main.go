package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func runLoop(vid *video.Video, m *menu.Menu) {
	var currTime, prevTime time.Time
	for !vid.Window.ShouldClose() {
		currTime = time.Now()
		dt := float32(currTime.Sub(prevTime)) / 1000000000
		glfw.PollEvents()
		m.ProcessHotkeys()
		ntf.Process(dt)
		vid.ResizeViewport()
		if !state.Global.MenuActive {
			if state.Global.CoreRunning {
				state.Global.Core.Run()
				if state.Global.Core.FrameTimeCallback != nil {
					state.Global.Core.FrameTimeCallback.Callback(state.Global.Core.FrameTimeCallback.Reference)
				}
				if state.Global.Core.AudioCallback != nil {
					state.Global.Core.AudioCallback.Callback()
				}
			}
			vid.Render()
		} else {
			input.Poll()
			m.Update(dt)
			vid.Render()
			m.Render(dt)
		}
		m.RenderNotifications()
		if state.Global.FastForward {
			glfw.SwapInterval(0)
		} else {
			glfw.SwapInterval(1)
		}
		vid.Window.SwapBuffers()
		prevTime = currTime
	}
}

func main() {
	err := settings.Load()
	if err != nil {
		log.Println("[Settings]: Loading failed:", err)
		log.Println("[Settings]: Using default settings")
	}

	var GLVersion string
	flag.StringVar(&state.Global.CorePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&state.Global.Verbose, "v", false, "Verbose logs")
	flag.StringVar(&GLVersion, "glver", settings.Defaults.GLVersion, "OpenGL version, possible values are 2.0, 2.1, 3.0, 3.1, 3.2, 4.1, 4.2")
	flag.BoolVar(&state.Global.LudOS, "ludos", false, "Expose the features related to LudOS")
	flag.Parse()
	args := flag.Args()

	if GLVersion != settings.Defaults.GLVersion {
		settings.Current.GLVersion = GLVersion
		if err := settings.Save(); err != nil {
			log.Fatalln("Failed to save settings:", err)
		}
	}

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	if !glfw.Init() {
		log.Fatalln("Failed to initialize glfw")
	}
	defer glfw.Terminate()

	state.Global.DB, err = scanner.LoadDB(settings.Current.DatabaseDirectory)
	if err != nil {
		log.Println("Can't load game database:", err)
	}

	playlists.Load()

	history.Load()

	vid := video.Init(settings.Current.VideoFullscreen, settings.Current.GLVersion)

	audio.Init()

	m := menu.Init(vid)

	core.Init(vid)

	input.Init(vid)

	if len(state.Global.CorePath) > 0 {
		err := core.Load(state.Global.CorePath)
		if err != nil {
			panic(err)
		}
	}

	if len(gamePath) > 0 {
		err := core.LoadGame(gamePath)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		} else {
			m.WarpToQuickMenu()
		}
	}

	// No game running? display the menu
	state.Global.MenuActive = !state.Global.CoreRunning

	runLoop(vid, m)

	// Unload and deinit in the core.
	core.Unload()
}
