package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"sync"
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
#include "libretro.h"
#cgo LDFLAGS: -ldl

void coreLog_cgo(enum retro_log_level level, const char *msg);
*/
import "C"

var mu sync.Mutex

//export coreLog
func coreLog(level C.enum_retro_log_level, msg *C.char) {
	fmt.Print("[Log]: ", C.GoString(msg))
}

//export coreEnvironment
func coreEnvironment(cmd C.unsigned, data unsafe.Pointer) bool {
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
	case C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE:
		cb := (*C.struct_retro_log_callback)(data)
		cb.log = (C.retro_log_printf_t)(C.coreLog_cgo)
		break
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

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func coreLoad(sofile string) {

	retroLoad(sofile)

	retroInit()

	v := retroAPIVersion()
	fmt.Println("Libretro API version:", v)
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

	gi := C.struct_retro_game_info{
		path: C.CString(filename),
		size: C.size_t(size),
	}

	si := retroGetSystemInfo()

	var libName = C.GoString(si.library_name)
	fmt.Println("  library_name:", libName)
	fmt.Println("  library_version:", C.GoString(si.library_version))
	fmt.Println("  valid_extensions:", C.GoString(si.valid_extensions))
	fmt.Println("  need_fullpath:", si.need_fullpath)
	fmt.Println("  block_extract:", si.block_extract)

	if !si.need_fullpath {
		bytes, err := slurp(filename, size)
		if err != nil {
			panic(err)
		}
		cstr := C.CString(string(bytes))
		gi.data = unsafe.Pointer(cstr)

	}

	ok := retroLoadGame(gi)
	if !ok {
		log.Fatal("The core failed to load the content.")
	}

	avi := retroGetSystemAVInfo()

	geom := retroGameGeometry{
		aspectRatio: float64(avi.geometry.aspect_ratio),
		baseWidth:   int(avi.geometry.base_width),
		baseHeight:  int(avi.geometry.base_height),
	}

	videoConfigure(geom)
	// Append the library name to the window title.
	if len(libName) > 0 {
		window.SetTitle("nanoarch - " + libName)
	}
	audioInit(int32(avi.timing.sample_rate))
}

//export coreVideoRefresh
func coreVideoRefresh(data unsafe.Pointer, width C.unsigned, height C.unsigned, pitch C.size_t) {
	videoRefresh(data, int32(width), int32(height), int32(pitch))
}

//export coreAudioSample
func coreAudioSample(left C.int16_t, right C.int16_t) {
	buf := []byte{byte(left), byte(right)}
	audioWrite(buf, 4)
}

//export coreAudioSampleBatch
func coreAudioSampleBatch(buf unsafe.Pointer, frames C.size_t) C.size_t {
	return C.size_t(audioWrite(C.GoBytes(buf, C.int(bufSize)), int32(frames)*4))
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

		retroRun()

		videoRender()

		window.SwapBuffers()
	}

	// Unload and deinit in the core.
	retroUnloadGame()
	retroDeinit()
}
