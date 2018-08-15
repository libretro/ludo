package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/input"
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/notifications"
)

// global state
var g struct {
	core        libretro.Core
	frameTimeCb libretro.FrameTimeCallback
	audioCb     libretro.AudioCallback
	coreRunning bool
	menuActive  bool
	verbose     bool
	corePath    string
	gamePath    string
	options     *Options
}

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
	for !window.ShouldClose() {
		glfw.SwapInterval(1)
		glfw.PollEvents()
		notifications.Process()
		if !g.menuActive {
			if g.coreRunning {
				g.core.Run()
				if g.frameTimeCb.Callback != nil {
					g.frameTimeCb.Callback(g.frameTimeCb.Reference)
				}
				if g.audioCb.Callback != nil {
					g.audioCb.Callback()
				}
			}
			videoRender()
		} else {
			input.Poll()
			menuInput()
			videoRender()
			menuRender()
		}
		renderNotifications()
		window.SwapBuffers()
	}
}

func main() {
	flag.StringVar(&g.corePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&g.verbose, "v", false, "Verbose logs")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	err := loadSettings()
	if err != nil {
		log.Println("[Settings]: Loading failed:", err)
		log.Println("[Settings]: Using default settings")
		saveSettings()
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	if len(g.corePath) > 0 {
		coreLoad(g.corePath)
	}

	video.winWidth = 320 * 3
	video.winHeight = 180 * 3

	videoConfigure(settings.VideoFullscreen)

	input.Init(window)

	if len(gamePath) > 0 {
		coreLoadGame(gamePath)
	}

	menuInit(window)

	// No game running? display the menu
	g.menuActive = !g.coreRunning

	runLoop()

	// Unload and deinit in the core.
	if g.coreRunning {
		g.core.UnloadGame()
		g.core.Deinit()
	}
}
