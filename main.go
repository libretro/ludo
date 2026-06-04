package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
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

func fulfillPendingScreenshot(vid *video.Video) bool {
	if state.PendingScreenshotName == "" {
		return false
	}

	name := state.PendingScreenshotName
	done := state.PendingScreenshotDone
	state.PendingScreenshotName = ""
	state.PendingScreenshotDone = nil

	err := vid.TakeScreenshotInFrontendFrame(name)
	if done != nil {
		done(err)
	}
	return true
}

func runLoop(vid *video.Video, m *menu.Menu) {
	var currTime time.Time
	prevTime := time.Now()
	for !vid.Window.ShouldClose() {
		currTime = time.Now()
		dt := float32(currTime.Sub(prevTime)) / 1000000000
		if state.MenuContextResetNeeded {
			m.ContextReset()
			state.MenuContextResetNeeded = false
		}
		glfw.PollEvents()
		m.ProcessHotkeys()
		ntf.Process(dt)
		m.UpdatePalette()
		input.Poll()
		if !state.MenuActive {
			if state.CoreRunning {
				vid.BeginCoreFrame()
				if !state.CoreSetSharedContext {
					vid.PrepareCoreContext()
				}
				state.Core.Run()
				vid.EndCoreFrame()
				if state.Core.FrameTimeCallback != nil {
					state.Core.FrameTimeCallback.Callback(state.Core.FrameTimeCallback.Reference)
				}
				if state.Core.AudioCallback != nil {
					state.Core.AudioCallback.Callback()
				}
			}
			coreGLState := vid.BeginFrontendFrame()
			vid.Render()
			vid.CachePresentedFrame()
			m.RenderNotifications()
			vid.EndFrontendFrame(coreGLState)
			frame++
			if frame%600 == 0 { // save sram about every 10 sec
				savefiles.SaveSRAM()
			}
		} else {
			coreGLState := vid.BeginFrontendFrame()
			m.Update(dt)
			vid.RenderPresentedFrame()
			if fulfillPendingScreenshot(vid) {
				vid.PrepareFrontendFrame()
				vid.RenderPresentedFrame()
			}
			m.Render(dt)
			m.RenderNotifications()
			vid.EndFrontendFrame(coreGLState)
		}
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

	// Match the normal menu-driven load path more closely: by the time a core is
	// loaded from the UI, GLFW has already pumped at least one event cycle on the
	// visible window. Some HW-core context recreation paths behave differently if
	// we load immediately from CLI before the first poll.
	glfw.PollEvents()

	if len(state.CorePath) > 0 {
		corePath := state.CorePath
		m.SetStartupAction(func() {
			if err := core.Load(corePath); err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				return
			}
			if len(gamePath) > 0 {
				m.SetStartupAction(func() {
					if err := core.LoadGame(gamePath); err != nil {
						ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
						return
					}
					m.WarpToQuickMenu()
					state.MenuActive = false
				})
			}
		})
	}

	// No game running? display the menu
	state.MenuActive = !state.CoreRunning

	runLoop(vid, m)

	// Unload and deinit in the core.
	core.Unload()
}
