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
#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>
#include <string.h>

void bridge_retro_deinit(void *f);
unsigned bridge_retro_api_version(void *f);
void bridge_retro_get_system_info(void *f, struct retro_system_info *si);
void bridge_retro_get_system_av_info(void *f, struct retro_system_av_info *si);
bool bridge_retro_set_environment(void *f, void *callback);
void bridge_retro_set_video_refresh(void *f, void *callback);
void bridge_retro_set_input_poll(void *f, void *callback);
void bridge_retro_set_input_state(void *f, void *callback);
void bridge_retro_set_audio_sample(void *f, void *callback);
void bridge_retro_set_audio_sample_batch(void *f, void *callback);
bool bridge_retro_load_game(void *f, struct retro_game_info *gi);
void bridge_retro_unload_game(void *f);
void bridge_retro_run(void *f);

bool coreEnvironment_cgo(unsigned cmd, void *data);
void coreVideoRefresh_cgo(void *data, unsigned width, unsigned height, size_t pitch);
void coreInputPoll_cgo();
void coreAudioSample_cgo(int16_t left, int16_t right);
size_t coreAudioSampleBatch_cgo(const int16_t *data, size_t frames);
int16_t coreInputState_cgo(unsigned port, unsigned device, unsigned index, unsigned id);
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

	mu.Lock()
	h := C.dlopen(C.CString(sofile), C.RTLD_NOW)
	if h == nil {
		log.Fatalf("error loading %s\n", sofile)
	}

	symRetroInit = C.dlsym(h, C.CString("retro_init"))
	symRetroDeinit = C.dlsym(h, C.CString("retro_deinit"))
	symRetroAPIVersion = C.dlsym(h, C.CString("retro_api_version"))
	symRetroGetSystemInfo = C.dlsym(h, C.CString("retro_get_system_info"))
	symRetroGetSystemAVInfo = C.dlsym(h, C.CString("retro_get_system_av_info"))
	symRetroSetEnvironment = C.dlsym(h, C.CString("retro_set_environment"))
	symRetroSetVideoRefresh = C.dlsym(h, C.CString("retro_set_video_refresh"))
	symRetroSetInputPoll = C.dlsym(h, C.CString("retro_set_input_poll"))
	symRetroSetInputState = C.dlsym(h, C.CString("retro_set_input_state"))
	symRetroSetAudioSample = C.dlsym(h, C.CString("retro_set_audio_sample"))
	symRetroSetAudioSampleBatch = C.dlsym(h, C.CString("retro_set_audio_sample_batch"))
	symRetroRun = C.dlsym(h, C.CString("retro_run"))
	symRetroLoadGame = C.dlsym(h, C.CString("retro_load_game"))
	symRetroUnloadGame = C.dlsym(h, C.CString("retro_unload_game"))
	mu.Unlock()

	C.bridge_retro_set_environment(symRetroSetEnvironment, C.coreEnvironment_cgo)
	C.bridge_retro_set_video_refresh(symRetroSetVideoRefresh, C.coreVideoRefresh_cgo)
	C.bridge_retro_set_input_poll(symRetroSetInputPoll, C.coreInputPoll_cgo)
	C.bridge_retro_set_input_state(symRetroSetInputState, C.coreInputState_cgo)
	C.bridge_retro_set_audio_sample(symRetroSetAudioSample, C.coreAudioSample_cgo)
	C.bridge_retro_set_audio_sample_batch(symRetroSetAudioSampleBatch, C.coreAudioSampleBatch_cgo)

	retroInit()

	v := C.bridge_retro_api_version(symRetroAPIVersion)
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

	si := C.struct_retro_system_info{}

	C.bridge_retro_get_system_info(symRetroGetSystemInfo, &si)

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

	ok := C.bridge_retro_load_game(symRetroLoadGame, &gi)
	if !ok {
		log.Fatal("The core failed to load the content.")
	}

	avi := C.struct_retro_system_av_info{}

	C.bridge_retro_get_system_av_info(symRetroGetSystemAVInfo, &avi)

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

		C.bridge_retro_run(symRetroRun)

		videoRender()

		window.SwapBuffers()
	}

	// Unload and deinit in the core.
	C.bridge_retro_unload_game(symRetroUnloadGame)
	C.bridge_retro_deinit(symRetroDeinit)
}
