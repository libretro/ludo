package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
	"github.com/libretro/ludo/netplay"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/savefiles"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var frame = 0

func runLoop(vid *video.Video, m *menu.Menu) {
	currTime := time.Now()
	prevTime := time.Now()
	for !vid.Window.ShouldClose() {
		currTime = time.Now()
		dt := float32(currTime.Sub(prevTime)) / 1000000000
		glfw.PollEvents()
		m.ProcessHotkeys()
		ntf.Process(dt)
		vid.ResizeViewport()
		m.UpdatePalette()

		state.ForcePause = vid.Window.GetKey(glfw.KeySpace) == glfw.Press

		if !state.MenuActive {
			if state.CoreRunning {
				netplay.Update()
			}
			vid.Render()
			frame++
			if frame%600 == 0 { // save sram about every 10 sec
				savefiles.SaveSRAM()
			}
		} else {
			input.Poll()
			m.Update(dt)
			vid.Render()
			m.Render(dt)
		}

		m.RenderNotifications()
		if state.FastForward {
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

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.StringVar(&state.CorePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&state.Verbose, "v", false, "Verbose logs")
	flag.BoolVar(&state.LudOS, "ludos", false, "Expose the features related to LudOS")
	flag.BoolVar(&netplay.Listen, "listen", false, "For the netplay server")
	flag.BoolVar(&netplay.Join, "join", false, "For the netplay client")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw", err)
	}
	defer glfw.Terminate()

	state.DB, err = scanner.LoadDB(settings.Current.DatabaseDirectory)
	if err != nil {
		log.Println("Can't load game database:", err)
	}

	playlists.Load()

	history.Load()

	vid := video.Init(settings.Current.VideoFullscreen)

	audio.Init()

	m := menu.Init(vid)

	core.Init(vid)

	input.Init(vid)

	if len(state.CorePath) > 0 {
		err := core.Load(state.CorePath)
		if err == nil {
			if len(gamePath) > 0 {
				err := core.LoadGame(gamePath)
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				} else {
					m.WarpToQuickMenu()
				}
			}
		} else {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		}
	}

	// No game running? display the menu
	state.MenuActive = !state.CoreRunning

	runLoop(vid, m)

	// Unload and deinit in the core.
	core.Unload()
}
