package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kivutar/go-playthemall/libretro"
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

	video.geom = libretro.GameGeometry{
		AspectRatio: 16.0 / 9.0,
		BaseWidth:   320 * 3,
		BaseHeight:  180 * 3,
	}

	videoConfigure(settings.VideoFullscreen)

	menuInit()

	// No game running? display the menu
	g.menuActive = !g.coreRunning

	if len(gamePath) > 0 {
		coreLoadGame(gamePath)
	}

	for !window.ShouldClose() {
		glfw.SwapInterval(1)
		glfw.PollEvents()
		processNotifications()
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
			inputPoll()
			menuInput()
			videoRender()
			menuRender()
		}
		renderNotifications()
		window.SwapBuffers()
	}

	// Unload and deinit in the core.
	if g.coreRunning {
		g.core.UnloadGame()
		g.core.Deinit()
	}
}
