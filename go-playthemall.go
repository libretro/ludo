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

var core libretro.Core

var menuActive bool

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
	core, _ = libretro.Load(sofile)
	core.SetEnvironment(environment)
	core.SetVideoRefresh(videoRefresh)
	core.SetInputPoll(inputPoll)
	core.SetInputState(inputState)
	core.SetAudioSample(audioSample)
	core.SetAudioSampleBatch(audioSampleBatch)
	core.Init()
	fmt.Println("Libretro API version:", core.APIVersion())
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

	si := core.GetSystemInfo()

	fmt.Println("  library_name:", si.LibraryName)
	fmt.Println("  library_version:", si.LibraryVersion)
	fmt.Println("  valid_extensions:", si.ValidExtensions)
	fmt.Println("  need_fullpath:", si.NeedFullpath)
	fmt.Println("  block_extract:", si.BlockExtract)

	if !si.NeedFullpath {
		bytes, err := slurp(filename, size)
		if err != nil {
			panic(err)
		}
		gi.SetData(bytes)
	}

	ok := core.LoadGame(gi)
	if !ok {
		log.Fatal("The core failed to load the content.")
	}

	avi := core.GetSystemAVInfo()

	// Create the video window, not-fullscreen.
	videoConfigure(avi.Geometry, false)

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		window.SetTitle("playthemall - " + si.LibraryName)
	}
	inputInit()
	audioInit(int32(avi.Timing.SampleRate))
}

func main() {
	var corePath = flag.String("L", "", "Path to the libretro core")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}
	if (len(*corePath) == 0 || len(gamePath) == 0) {
		log.Fatalln("Usage: go-playthemall -L <core> <game>")
		return
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	coreLoad(*corePath)
	coreLoadGame(gamePath)
	menuInit()

	for !window.ShouldClose() {
		glfw.PollEvents()
		if !menuActive {
			core.Run()
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
	core.UnloadGame()
	core.Deinit()
}
