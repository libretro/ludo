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
	"unsafe"
)

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

var symRetroInit unsafe.Pointer
var symRetroDeinit unsafe.Pointer
var symRetroAPIVersion unsafe.Pointer
var symRetroGetSystemInfo unsafe.Pointer
var symRetroGetSystemAVInfo unsafe.Pointer
var symRetroSetEnvironment unsafe.Pointer
var symRetroSetVideoRefresh unsafe.Pointer
var symRetroSetInputPoll unsafe.Pointer
var symRetroSetInputState unsafe.Pointer
var symRetroSetAudioSample unsafe.Pointer
var symRetroSetAudioSampleBatch unsafe.Pointer
var symRetroRun unsafe.Pointer
var symRetroLoadGame unsafe.Pointer
var symRetroUnloadGame unsafe.Pointer

func retroLoad(sofile string) {
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
}

func retroInit() {
	C.bridge_retro_init(symRetroInit)
}

func retroAPIVersion() uint {
	return uint(C.bridge_retro_api_version(symRetroAPIVersion))
}

func retroDeinit() {
	C.bridge_retro_deinit(symRetroDeinit)
}

func retroRun() {
	C.bridge_retro_run(symRetroRun)
}

func retroGetSystemInfo() retroSystemInfo {
	rsi := C.struct_retro_system_info{}
	C.bridge_retro_get_system_info(symRetroGetSystemInfo, &rsi)
	return retroSystemInfo{
		libraryName:     C.GoString(rsi.library_name),
		libraryVersion:  C.GoString(rsi.library_version),
		validExtensions: C.GoString(rsi.valid_extensions),
		needFullpath:    bool(rsi.need_fullpath),
		blockExtract:    bool(rsi.block_extract),
	}
}

func retroGetSystemAVInfo() C.struct_retro_system_av_info {
	si := C.struct_retro_system_av_info{}
	C.bridge_retro_get_system_av_info(symRetroGetSystemAVInfo, &si)
	return si
}

func retroLoadGame(gi retroGameInfo) bool {
	rgi := C.struct_retro_game_info{}
	rgi.path = C.CString(gi.path)
	rgi.size = C.size_t(gi.size)
	rgi.data = gi.data
	return bool(C.bridge_retro_load_game(symRetroLoadGame, &rgi))
}

func retroUnloadGame() {
	C.bridge_retro_unload_game(symRetroUnloadGame)
}
