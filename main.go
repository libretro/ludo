package main

import (
	"flag"
	"fmt"
	"log"
	"math"
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

const (
	defaultCoreFPS     = 60.0
	maxLoopDelta       = 250 * time.Millisecond
	maxPacingLagFrames = 2
	maxAutoSwapInt     = 4
	pacingMatchSkew    = 0.02
)

func coreFrameDuration() time.Duration {
	fps := state.CoreFPS
	if fps <= 0 {
		fps = defaultCoreFPS
	}

	return time.Duration(float64(time.Second) / fps)
}

func frameTimeUsec(frameDuration time.Duration) int64 {
	if state.Core != nil && state.Core.FrameTimeCallback != nil {
		if state.FastForward && state.Core.FrameTimeCallback.Reference > 0 {
			return state.Core.FrameTimeCallback.Reference
		}

		usec := frameDuration / time.Microsecond
		if usec > 0 {
			return int64(usec)
		}

		if state.Core.FrameTimeCallback.Reference > 0 {
			return state.Core.FrameTimeCallback.Reference
		}
	}

	return int64(coreFrameDuration() / time.Microsecond)
}

func runCoreFrame(frameDuration time.Duration) {
	if state.Core == nil {
		return
	}

	if state.Core.FrameTimeCallback != nil {
		state.Core.FrameTimeCallback.Callback(frameTimeUsec(frameDuration))
	}

	state.Core.Run()
}

func updateSwapInterval(current *int, target int) {
	if *current == target {
		return
	}

	glfw.SwapInterval(target)
	*current = target
}

func relativeSkew(a, b float64) float64 {
	if a <= 0 || b <= 0 {
		return math.Inf(1)
	}
	return math.Abs(a-b) / b
}

func displayRefreshHz(vid *video.Video) float64 {
	if vid == nil || vid.Window == nil {
		return 0
	}

	monitor := vid.Window.GetMonitor()
	if monitor == nil {
		monitors := glfw.GetMonitors()
		if len(monitors) > 0 {
			index := settings.Current.VideoMonitorIndex
			if index < 0 || index >= len(monitors) {
				index = 0
			}
			monitor = monitors[index]
		}
	}

	if monitor == nil {
		return 0
	}

	mode := monitor.GetVideoMode()
	if mode == nil || mode.RefreshRate <= 0 {
		return 0
	}

	return float64(mode.RefreshRate)
}

func autoSwapInterval(refreshHz, coreFPS float64) int {
	if refreshHz <= 0 || coreFPS <= 0 {
		return 1
	}

	bestInterval := 1
	bestSkew := relativeSkew(refreshHz, coreFPS)

	for interval := 2; interval <= maxAutoSwapInt; interval++ {
		effectiveRefresh := refreshHz / float64(interval)
		skew := relativeSkew(effectiveRefresh, coreFPS)
		if skew <= pacingMatchSkew && skew < bestSkew {
			bestInterval = interval
			bestSkew = skew
		}
	}

	return bestInterval
}

func effectiveRefreshHz(refreshHz float64, swapInterval int) float64 {
	if refreshHz <= 0 {
		return 0
	}
	if swapInterval <= 0 {
		return refreshHz
	}
	return refreshHz / float64(swapInterval)
}

func blocksOnSwap(refreshHz, coreFPS float64) bool {
	return relativeSkew(refreshHz, coreFPS) <= pacingMatchSkew
}

func runLoop(vid *video.Video, m *menu.Menu) {
	var currTime time.Time
	prevTime := time.Now()
	accumulator := time.Duration(0)
	swapInterval := -1
	pacingPrimed := false
	for !vid.Window.ShouldClose() {
		currTime = time.Now()
		loopDelta := currTime.Sub(prevTime)
		prevTime = currTime
		if loopDelta > maxLoopDelta {
			loopDelta = maxLoopDelta
		}

		dt := float32(loopDelta) / 1000000000
		glfw.PollEvents()
		m.ProcessHotkeys()
		ntf.Process(dt)
		vid.ResizeViewport()
		m.UpdatePalette()
		input.Poll()
		if !state.MenuActive {
			ranFrames := 0
			if state.CoreRunning {
				if state.FastForward {
					updateSwapInterval(&swapInterval, 0)
					audio.SetVideoTiming(0, 0)
					runCoreFrame(coreFrameDuration())
					ranFrames = 1
					accumulator = 0
					pacingPrimed = false
				} else {
					frameDuration := coreFrameDuration()
					refreshHz := displayRefreshHz(vid)
					targetSwap := autoSwapInterval(refreshHz, state.CoreFPS)
					updateSwapInterval(&swapInterval, targetSwap)
					effectiveRefresh := effectiveRefreshHz(refreshHz, targetSwap)
					audio.SetVideoTiming(refreshHz, targetSwap)

					if blocksOnSwap(effectiveRefresh, state.CoreFPS) {
						runCoreFrame(frameDuration)
						ranFrames = 1
						accumulator = 0
						pacingPrimed = false
					} else {
						if !pacingPrimed {
							accumulator = frameDuration
							pacingPrimed = true
						}

						accumulator += loopDelta
						maxAccumulator := frameDuration * maxPacingLagFrames
						if accumulator > maxAccumulator {
							accumulator = maxAccumulator
						}

						if accumulator >= frameDuration {
							runCoreFrame(frameDuration)
							accumulator -= frameDuration
							if accumulator > frameDuration {
								accumulator = frameDuration
							}
							ranFrames = 1
						}
					}
				}

				if state.Core.AudioCallback != nil && ranFrames > 0 {
					state.Core.AudioCallback.Callback()
				}
			} else {
				updateSwapInterval(&swapInterval, 1)
				audio.SetVideoTiming(0, 0)
				accumulator = 0
				pacingPrimed = false
			}

			vid.Render()
			for i := 0; i < ranFrames; i++ {
				frame++
				if frame%600 == 0 { // save sram about every 10 sec
					savefiles.SaveSRAM()
				}
			}
		} else {
			updateSwapInterval(&swapInterval, 1)
			audio.SetVideoTiming(0, 0)
			accumulator = 0
			pacingPrimed = false
			m.Update(dt)
			vid.Render()
			m.Render(dt)
		}
		m.RenderNotifications()
		vid.Window.SwapBuffers()
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
