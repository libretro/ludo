package main

import (
	"flag"
	"fmt"
	"libretro"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// global state
var g struct {
	core        libretro.Core
	frameTimeCb libretro.FrameTimeCallback
	audioCb     libretro.AudioCallback
	coreRunning bool
	menuActive  bool
	gamePath    string
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var logLevels = map[uint32]string{
	libretro.LogLevelDebug: "DEBUG",
	libretro.LogLevelInfo:  "INFO",
	libretro.LogLevelWarn:  "WARN",
	libretro.LogLevelError: "ERROR",
	libretro.LogLevelDummy: "DUMMY",
}

func nanoLog(level uint32, str string) {
	fmt.Printf("[%s]: %s", logLevels[level], str)
}

func coreLoad(sofile string) {
	if g.coreRunning {
		g.core.UnloadGame()
		g.core.Deinit()
		g.gamePath = ""
	}

	g.core, _ = libretro.Load(sofile)
	g.core.SetEnvironment(environment)
	g.core.SetVideoRefresh(videoRefresh)
	g.core.SetInputPoll(inputPoll)
	g.core.SetInputState(inputState)
	g.core.SetAudioSample(audioSample)
	g.core.SetAudioSampleBatch(audioSampleBatch)
	g.core.Init()

	// Append the library name to the window title.
	si := g.core.GetSystemInfo()
	if len(si.LibraryName) > 0 {
		if window != nil {
			window.SetTitle("Play Them All - " + si.LibraryName)
		}
		fmt.Println("[Libretro]: Name:", si.LibraryName)
		fmt.Println("[Libretro]: Version:", si.LibraryVersion)
		fmt.Println("[Libretro]: Valid extensions:", si.ValidExtensions)
		fmt.Println("[Libretro]: Need fullpath:", si.NeedFullpath)
		fmt.Println("[Libretro]: Block extract:", si.BlockExtract)
	}

	notify("Core loaded: "+si.LibraryName, 240)
	fmt.Println("[Libretro]: API version:", g.core.APIVersion())
}

func coreLoadGame(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}

	size := fi.Size()

	fmt.Println("[Libretro]: ROM size:", size)

	gi := libretro.GameInfo{
		Path: filename,
		Size: size,
	}

	si := g.core.GetSystemInfo()

	if !si.NeedFullpath {
		bytes, err := slurp(filename)
		if err != nil {
			panic(err)
		}
		gi.SetData(bytes)
	}

	ok := g.core.LoadGame(gi)
	if !ok {
		notify("The core failed to load the content.", 240)
		fmt.Println("[Libretro]: The core failed to load the content.")
		g.coreRunning = false
		return
	}

	avi := g.core.GetSystemAVInfo()

	// Create the video window
	videoConfigure(avi.Geometry, settings.VideoFullscreen)

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		window.SetTitle("Play Them All - " + si.LibraryName)
	}

	inputInit()
	audioInit(int32(avi.Timing.SampleRate))
	if g.audioCb.SetState != nil {
		g.audioCb.SetState(true)
	}

	g.coreRunning = true
	g.menuActive = false
	g.gamePath = filename
	menuInit()
	notify("Game loaded: "+filename, 240)
}

func main() {
	var corePath = flag.String("L", "", "Path to the libretro core")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	err := loadSettings()
	if err != nil {
		fmt.Println("[Settings]: Loading failed:", err)
		fmt.Println("[Settings]: Using default settings")
		saveSettings()
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	if len(*corePath) > 0 {
		coreLoad(*corePath)
	}

	if len(gamePath) > 0 {
		coreLoadGame(gamePath)
	}

	// No game running? display the menu with a dummy geometry
	if !g.coreRunning {
		geom := libretro.GameGeometry{
			AspectRatio: 320.0 / 240.0,
			BaseWidth:   320,
			BaseHeight:  240,
		}
		videoConfigure(geom, settings.VideoFullscreen)
		menuInit()
		g.menuActive = true
	}

	for !window.ShouldClose() {
		glfw.PollEvents()
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
			renderMenuList()
		}
		window.SwapBuffers()
	}

	// Unload and deinit in the core.
	if g.coreRunning {
		g.core.UnloadGame()
		g.core.Deinit()
	}
}
