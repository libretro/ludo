package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/dat"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/savefiles"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
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
		w, h := vid.Window.GetFramebufferSize()
		vid.Font.UpdateResolution(w, h)
		vid.BoldFont.UpdateResolution(w, h)
		m.UpdatePalette()
		input.Poll()
		if !state.MenuActive {
			if state.CoreRunning {
				state.Core.Run()
				if state.Core.FrameTimeCallback != nil {
					state.Core.FrameTimeCallback.Callback(state.Core.FrameTimeCallback.Reference)
				}
				if state.Core.AudioCallback != nil {
					state.Core.AudioCallback.Callback()
				}
			}
			vid.Render()
			frame++
			if frame%600 == 0 { // save sram about every 10 sec
				savefiles.SaveSRAM()
			}
		} else {
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

	// ExitOnError causes flags to quit after displaying help.
	// (--help counts as an error)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// customize help message
	flag.CommandLine.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] [content]\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	// set arguments
	flag.StringVar(&state.CorePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&state.Verbose, "v", false, "Verbose logs")
	flag.BoolVar(&state.LudOS, "ludos", false, "Expose the features related to LudOS")
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
					scanner.ScanFile(gamePath, func(game dat.Game) {
						name := game.Name
						if name == "" {
							name = utils.FileName(gamePath)
						}
						history.Push(history.Game{
							Path:     gamePath,
							Name:     name,
							System:   game.System,
							CorePath: state.CorePath,
						})
						history.Load()
						m.WarpToQuickMenu()
					})
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
