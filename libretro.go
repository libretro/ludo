package main

/*
#include "libretro.h"

void bridge_retro_init(void *f);
void bridge_retro_run(void *f);
*/
import "C"
import "unsafe"

type retroGameGeometry struct {
	aspectRatio float64
	baseWidth   int
	baseHeight  int
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

func retroInit() {
	C.bridge_retro_init(symRetroInit)
}

func retroRun() {
	C.bridge_retro_run(symRetroRun)
}
