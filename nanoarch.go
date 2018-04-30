package main

import (
	"C"
	"flag"
	"fmt"
	"libretro"
	"log"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var core libretro.Core

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func coreLoad(sofile string) {
	core = libretro.Load(sofile)
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
		cstr := C.CString(string(bytes))
		gi.Data = unsafe.Pointer(cstr)
	}

	ok := core.LoadGame(gi)
	if !ok {
		log.Fatal("The core failed to load the content.")
	}

	avi := core.GetSystemAVInfo()

	videoConfigure(avi.Geometry)
	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		window.SetTitle("nanoarch - " + si.LibraryName)
	}
	audioInit(int32(avi.Timing.SampleRate))
}

func main() {
	var corePath = flag.String("L", "", "Path to the libretro core")
	var gamePath = flag.String("G", "", "Path to the game")
	flag.Parse()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	coreLoad(*corePath)
	coreLoadGame(*gamePath)

	for !window.ShouldClose() {
		glfw.PollEvents()
		core.Run()
		videoRender()
		window.SwapBuffers()
	}

	// Unload and deinit in the core.
	core.UnloadGame()
	core.Deinit()
}
