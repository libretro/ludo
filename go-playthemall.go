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
	coreRunning bool
	menuActive  bool
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
		window.SetTitle("playthemall - " + si.LibraryName)
		fmt.Println("[Libretro]: Name:", si.LibraryName)
		fmt.Println("[Libretro]: Version:", si.LibraryVersion)
		fmt.Println("[Libretro]: Valid extensions:", si.ValidExtensions)
		fmt.Println("[Libretro]: Need fullpath:", si.NeedFullpath)
		fmt.Println("[Libretro]: Block extract:", si.BlockExtract)
	}

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

	fmt.Println("ROM size:", size)

	gi := libretro.GameInfo{
		Path: filename,
		Size: size,
	}

	si := g.core.GetSystemInfo()

	if !si.NeedFullpath {
		bytes, err := slurp(filename, size)
		if err != nil {
			panic(err)
		}
		gi.SetData(bytes)
	}

	ok := g.core.LoadGame(gi)
	if !ok {
		log.Fatal("The core failed to load the content.")
	}

	avi := g.core.GetSystemAVInfo()

	// Create the video window, not-fullscreen.
	videoConfigure(avi.Geometry, false)

	inputInit()
	audioInit(int32(avi.Timing.SampleRate))

	g.coreRunning = true
}

func main() {
	var corePath = flag.String("L", "", "Path to the libretro core")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
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
		videoConfigure(geom, false)
		g.menuActive = true
	}

	menuInit()

	for !window.ShouldClose() {
		glfw.PollEvents()
		if !g.menuActive {
			g.core.Run()
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
