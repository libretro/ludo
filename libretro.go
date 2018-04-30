package main

/*
#include "libretro.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <dlfcn.h>

void bridge_retro_init(void *f);
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
import (
	"log"
	"sync"
	"unsafe"
)

type retro struct {
	handle                      unsafe.Pointer
	symRetroInit                unsafe.Pointer
	symRetroDeinit              unsafe.Pointer
	symRetroAPIVersion          unsafe.Pointer
	symRetroGetSystemInfo       unsafe.Pointer
	symRetroGetSystemAVInfo     unsafe.Pointer
	symRetroSetEnvironment      unsafe.Pointer
	symRetroSetVideoRefresh     unsafe.Pointer
	symRetroSetInputPoll        unsafe.Pointer
	symRetroSetInputState       unsafe.Pointer
	symRetroSetAudioSample      unsafe.Pointer
	symRetroSetAudioSampleBatch unsafe.Pointer
	symRetroRun                 unsafe.Pointer
	symRetroLoadGame            unsafe.Pointer
	symRetroUnloadGame          unsafe.Pointer
}

type retroGameGeometry struct {
	aspectRatio float64
	baseWidth   int
	baseHeight  int
}

type retroGameInfo struct {
	path string
	size int64
	data unsafe.Pointer
}

type retroSystemInfo struct {
	libraryName     string
	libraryVersion  string
	validExtensions string
	needFullpath    bool
	blockExtract    bool
}

var (
	retroPixelFormat0RGB1555 = uint32(C.RETRO_PIXEL_FORMAT_0RGB1555)
	retroPixelFormatXRGB8888 = uint32(C.RETRO_PIXEL_FORMAT_XRGB8888)
	retroPixelFormatRGB565   = uint32(C.RETRO_PIXEL_FORMAT_RGB565)
)

var (
	retroDeviceIDJoypadB      = uint32(C.RETRO_DEVICE_ID_JOYPAD_B)
	retroDeviceIDJoypadY      = uint32(C.RETRO_DEVICE_ID_JOYPAD_Y)
	retroDeviceIDJoypadSelect = uint32(C.RETRO_DEVICE_ID_JOYPAD_SELECT)
	retroDeviceIDJoypadStart  = uint32(C.RETRO_DEVICE_ID_JOYPAD_START)
	retroDeviceIDJoypadUp     = uint32(C.RETRO_DEVICE_ID_JOYPAD_UP)
	retroDeviceIDJoypadDown   = uint32(C.RETRO_DEVICE_ID_JOYPAD_DOWN)
	retroDeviceIDJoypadLeft   = uint32(C.RETRO_DEVICE_ID_JOYPAD_LEFT)
	retroDeviceIDJoypadRight  = uint32(C.RETRO_DEVICE_ID_JOYPAD_RIGHT)
	retroDeviceIDJoypadA      = uint32(C.RETRO_DEVICE_ID_JOYPAD_A)
	retroDeviceIDJoypadX      = uint32(C.RETRO_DEVICE_ID_JOYPAD_X)
	retroDeviceIDJoypadL      = uint32(C.RETRO_DEVICE_ID_JOYPAD_L)
	retroDeviceIDJoypadR      = uint32(C.RETRO_DEVICE_ID_JOYPAD_R)
	retroDeviceIDJoypadL2     = uint32(C.RETRO_DEVICE_ID_JOYPAD_L2)
	retroDeviceIDJoypadR2     = uint32(C.RETRO_DEVICE_ID_JOYPAD_R2)
	retroDeviceIDJoypadL3     = uint32(C.RETRO_DEVICE_ID_JOYPAD_L3)
	retroDeviceIDJoypadR3     = uint32(C.RETRO_DEVICE_ID_JOYPAD_R3)
)

var mu sync.Mutex

func retroLoad(sofile string) retro {
	r := retro{}

	mu.Lock()
	r.handle = C.dlopen(C.CString(sofile), C.RTLD_NOW)
	if r.handle == nil {
		log.Fatalf("error loading %s\n", sofile)
	}

	r.symRetroInit = C.dlsym(r.handle, C.CString("retro_init"))
	r.symRetroDeinit = C.dlsym(r.handle, C.CString("retro_deinit"))
	r.symRetroAPIVersion = C.dlsym(r.handle, C.CString("retro_api_version"))
	r.symRetroGetSystemInfo = C.dlsym(r.handle, C.CString("retro_get_system_info"))
	r.symRetroGetSystemAVInfo = C.dlsym(r.handle, C.CString("retro_get_system_av_info"))
	r.symRetroSetEnvironment = C.dlsym(r.handle, C.CString("retro_set_environment"))
	r.symRetroSetVideoRefresh = C.dlsym(r.handle, C.CString("retro_set_video_refresh"))
	r.symRetroSetInputPoll = C.dlsym(r.handle, C.CString("retro_set_input_poll"))
	r.symRetroSetInputState = C.dlsym(r.handle, C.CString("retro_set_input_state"))
	r.symRetroSetAudioSample = C.dlsym(r.handle, C.CString("retro_set_audio_sample"))
	r.symRetroSetAudioSampleBatch = C.dlsym(r.handle, C.CString("retro_set_audio_sample_batch"))
	r.symRetroRun = C.dlsym(r.handle, C.CString("retro_run"))
	r.symRetroLoadGame = C.dlsym(r.handle, C.CString("retro_load_game"))
	r.symRetroUnloadGame = C.dlsym(r.handle, C.CString("retro_unload_game"))
	mu.Unlock()

	C.bridge_retro_set_environment(r.symRetroSetEnvironment, C.coreEnvironment_cgo)
	C.bridge_retro_set_video_refresh(r.symRetroSetVideoRefresh, C.coreVideoRefresh_cgo)
	C.bridge_retro_set_input_poll(r.symRetroSetInputPoll, C.coreInputPoll_cgo)
	C.bridge_retro_set_input_state(r.symRetroSetInputState, C.coreInputState_cgo)
	C.bridge_retro_set_audio_sample(r.symRetroSetAudioSample, C.coreAudioSample_cgo)
	C.bridge_retro_set_audio_sample_batch(r.symRetroSetAudioSampleBatch, C.coreAudioSampleBatch_cgo)

	return r
}

func (r retro) Init() {
	C.bridge_retro_init(r.symRetroInit)
}

func (r retro) APIVersion() uint {
	return uint(C.bridge_retro_api_version(r.symRetroAPIVersion))
}

func (r retro) Deinit() {
	C.bridge_retro_deinit(r.symRetroDeinit)
}

func (r retro) Run() {
	C.bridge_retro_run(r.symRetroRun)
}

func (r retro) GetSystemInfo() retroSystemInfo {
	rsi := C.struct_retro_system_info{}
	C.bridge_retro_get_system_info(r.symRetroGetSystemInfo, &rsi)
	return retroSystemInfo{
		libraryName:     C.GoString(rsi.library_name),
		libraryVersion:  C.GoString(rsi.library_version),
		validExtensions: C.GoString(rsi.valid_extensions),
		needFullpath:    bool(rsi.need_fullpath),
		blockExtract:    bool(rsi.block_extract),
	}
}

func (r retro) GetSystemAVInfo() C.struct_retro_system_av_info {
	si := C.struct_retro_system_av_info{}
	C.bridge_retro_get_system_av_info(r.symRetroGetSystemAVInfo, &si)
	return si
}

func (r retro) LoadGame(gi retroGameInfo) bool {
	rgi := C.struct_retro_game_info{}
	rgi.path = C.CString(gi.path)
	rgi.size = C.size_t(gi.size)
	rgi.data = gi.data
	return bool(C.bridge_retro_load_game(r.symRetroLoadGame, &rgi))
}

func (r retro) UnloadGame() {
	C.bridge_retro_unload_game(r.symRetroUnloadGame)
}
