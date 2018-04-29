package main

/*
#include "libretro.h"
*/
import "C"

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
