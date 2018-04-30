package main

import (
	"flag"
	"fmt"
	"libretro"
	"log"
	"os"
	"os/user"
	"runtime"
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
#include "vendor/libretro/libretro.h"
#include <stdbool.h>
#include <stdarg.h>
#include <stdio.h>
#cgo LDFLAGS: -ldl
*/
import "C"

var core libretro.Core

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func environment(cmd uint, data unsafe.Pointer) bool {
	switch cmd {
	case C.RETRO_ENVIRONMENT_GET_USERNAME:
		username := (**C.char)(data)
		currentUser, err := user.Current()
		if err != nil {
			*username = C.CString("")
		} else {
			*username = C.CString(currentUser.Username)
		}
		break
	// case C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE:
	// 	cb := (*C.struct_retro_log_callback)(data)
	// 	cb.log = (C.retro_log_printf_t)(C.coreLog_cgo)
	// 	break
	case C.RETRO_ENVIRONMENT_GET_CAN_DUPE:
		bval := (*C.bool)(data)
		*bval = C.bool(true)
		break
	case C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT:
		format := (*C.enum_retro_pixel_format)(data)
		if *format > C.RETRO_PIXEL_FORMAT_RGB565 {
			return false
		}
		return videoSetPixelFormat(*format)
	case C.RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY:
	case C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY:
		path := (**C.char)(data)
		*path = C.CString(".")
		return true
	case C.RETRO_ENVIRONMENT_SHUTDOWN:
		window.SetShouldClose(true)
		return true
	case C.RETRO_ENVIRONMENT_GET_VARIABLE:
		variable := (*C.struct_retro_variable)(data)
		fmt.Println("[Env]: get variable:", C.GoString(variable.key))
		return false
	default:
		//fmt.Println("[Env]: command not implemented", cmd)
		return false
	}
	return true
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
