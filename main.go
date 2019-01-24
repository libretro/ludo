package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
	"github.com/libretro/ludo/notifications"
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

func runLoop(vid *video.Video) {
	for !vid.Window.ShouldClose() {
		glfw.PollEvents()
		notifications.Process()
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
			menu.Update()
			vid.Render()
			menu.Render()
		}
		input.ProcessActions()
		menu.RenderNotifications()
		glfw.SwapInterval(1)
		vid.Window.SwapBuffers()
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
	flag.Parse()
	args := flag.Args()

	if GLVersion != settings.Defaults.GLVersion {
		settings.Current.GLVersion = GLVersion
		settings.Save()
	}

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	state.Global.DB, err = scanner.LoadDB(settings.Current.DatabaseDirectory)
	if err != nil {
		log.Println("Can't load game database:", err)
	}

	playlists.LoadPlaylists()

	vid := video.Init(settings.Current.VideoFullscreen, settings.Current.GLVersion)

	m := menu.Init(vid)
	m.ContextReset()

	core.Init(vid, m)

	input.Init(vid, m)

	if len(state.Global.CorePath) > 0 {
		err := core.Load(state.Global.CorePath)
		if err != nil {
			panic(err)
		}
	}

	if len(gamePath) > 0 {
		err := core.LoadGame(gamePath)
		if err != nil {
			notifications.DisplayAndLog("error", "Menu", err.Error())
		} else {
			m.WarpToQuickMenu()
		}
	}

	// No game running? display the menu
	state.Global.MenuActive = !state.Global.CoreRunning

	runLoop(vid)

	// Unload and deinit in the core.
	core.Unload()
}
