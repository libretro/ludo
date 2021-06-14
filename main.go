package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
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

// Micromage, getPlayer Score
// playerIndex valid values
// 0-3 : Player 1 to Player 4
// 4 : Total Score
func getScore(playerIndex int) uint16 {
	systemRamPointer := state.Core.GetMemoryData(libretro.MemorySystemRAM)

	// Micromages
	// 0x5F5 (BYTE) is Player 1 (00-99)
	// 0x5F6 (BYTE) is Player 2 (00-99)
	// 0x5F7 (BYTE) is Player 3 (00-99)
	// 0x5F8 (BYTE) is Player 4 (00-99)
	// 0x5F9 (BYTE) contains the total score of the team

	// 0x5FA (BYTE) is Player 1 Overflow
	// 0x5FB (BYTE) is Player 2 Overflow
	// 0x5FC (BYTE) is Player 3 Overflow
	// 0x5FD (BYTE) is Player 4 Overflow
	// 0x5FE (BYTE) is total score overflow

	var bcdScore uint16 = (uint16)(*(*uint8)(unsafe.Pointer(uintptr(systemRamPointer) + uintptr(0x5FA+playerIndex))) << 8)
	bcdScore |= (uint16)(*(*uint8)(unsafe.Pointer(uintptr(systemRamPointer) + uintptr(0x5F5+playerIndex))))

	var actualScore = bcdScore & 0xF
	bcdScore >>= 4
	actualScore += (bcdScore & 0xF) * 10
	bcdScore >>= 4
	actualScore += (bcdScore & 0xF) * 100
	bcdScore >>= 4
	actualScore += (bcdScore & 0xF) * 1000

	actualScore *= 100

	return actualScore
}

// Micromages: Return the number of players playing
func getNumPlayers() uint8 {
	systemRamPointer := state.Core.GetMemoryData(libretro.MemorySystemRAM)
	numPlayers := *(*uint8)(unsafe.Pointer(uintptr(systemRamPointer) + uintptr(0x803)))

	return numPlayers + 1
}

func checkAchievements() {
	fmt.Printf("Player1: %d\n", getScore(0))
	fmt.Printf("Player2: %d\n", getScore(1))
	fmt.Printf("Player3: %d\n", getScore(2))
	fmt.Printf("Player4: %d\n", getScore(3))
	fmt.Printf("Total: %d\n", getScore(4))
	fmt.Printf("Numer of Players: %d\n", getNumPlayers())
}

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
			checkAchievements()

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
