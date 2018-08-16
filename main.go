package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/core"
	"github.com/libretro/go-playthemall/input"
	"github.com/libretro/go-playthemall/menu"
	"github.com/libretro/go-playthemall/notifications"
	"github.com/libretro/go-playthemall/options"
	"github.com/libretro/go-playthemall/settings"
	"github.com/libretro/go-playthemall/state"
	"github.com/libretro/go-playthemall/video"
)

var opts *options.Options
var vid *video.Video

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	// Create base folders
	usr, _ := user.Current()
	os.Mkdir(usr.HomeDir+"/.playthemall/", 0777)
	os.Mkdir(usr.HomeDir+"/.playthemall/playlists/", 0777)
	os.Mkdir(usr.HomeDir+"/.playthemall/savefiles/", 0777)
	os.Mkdir(usr.HomeDir+"/.playthemall/savestates/", 0777)
	os.Mkdir(usr.HomeDir+"/.playthemall/screenshots/", 0777)
	os.Mkdir(usr.HomeDir+"/.playthemall/system/", 0777)
}

func runLoop() {
	for !vid.Window.ShouldClose() {
		glfw.SwapInterval(1)
		glfw.PollEvents()
		notifications.Process()
		if !state.Global.MenuActive {
			if state.Global.CoreRunning {
				state.Global.Core.Run()
				if state.Global.FrameTimeCb.Callback != nil {
					state.Global.FrameTimeCb.Callback(state.Global.FrameTimeCb.Reference)
				}
				if state.Global.AudioCb.Callback != nil {
					state.Global.AudioCb.Callback()
				}
			}
			vid.Render()
		} else {
			input.Poll()
			menu.Input()
			vid.Render()
			menu.Render()
		}
		vid.RenderNotifications()
		vid.Window.SwapBuffers()
	}
}

func main() {
	flag.StringVar(&state.Global.CorePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&state.Global.Verbose, "v", false, "Verbose logs")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	err := settings.Load()
	if err != nil {
		log.Println("[Settings]: Loading failed:", err)
		log.Println("[Settings]: Using default settings")
		settings.Save()
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	if len(state.Global.CorePath) > 0 {
		core.Load(state.Global.CorePath)
	}

	vid = video.Init(settings.Settings.VideoFullscreen)
	menu.ContextReset()

	core.Init(vid, opts)

	input.Init(vid)

	if len(gamePath) > 0 {
		core.LoadGame(gamePath)
	}

	menu.Init(vid, opts)

	// No game running? display the menu
	state.Global.MenuActive = !state.Global.CoreRunning

	runLoop()

	// Unload and deinit in the core.
	if state.Global.CoreRunning {
		state.Global.Core.UnloadGame()
		state.Global.Core.Deinit()
	}
}
