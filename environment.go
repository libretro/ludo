package main

import (
	"fmt"
	"libretro"
	"os/user"
	"unsafe"
)

/*
#include "vendor/libretro/libretro.h"
*/
import "C"

func environment(cmd uint32, data unsafe.Pointer) bool {
	switch cmd {
	case libretro.EnvironmentGetUsername:
		username := (**C.char)(data)
		currentUser, err := user.Current()
		if err != nil {
			*username = C.CString("")
		} else {
			*username = C.CString(currentUser.Username)
		}
		break
	//case libretro.EnvironmentGetLogInterface:
	// 	cb := (*C.struct_retro_log_callback)(data)
	// 	cb.log = (C.retro_log_printf_t)(C.coreLog_cgo)
	// 	break
	case libretro.EnvironmentGetCanDupe:
		bval := (*C.bool)(data)
		*bval = C.bool(true)
		break
	case libretro.EnvironmentSetPixelFormat:
		format := (*C.enum_retro_pixel_format)(data)
		if *format > C.RETRO_PIXEL_FORMAT_RGB565 {
			return false
		}
		return videoSetPixelFormat(*format)
	case libretro.EnvironmentGetSystemDirectory:
	case libretro.EnvironmentGetSaveDirectory:
		path := (**C.char)(data)
		*path = C.CString(".")
		return true
	case libretro.EnvironmentShutdown:
		window.SetShouldClose(true)
		return true
	case libretro.EnvironmentGetVariable:
		variable := (*C.struct_retro_variable)(data)
		fmt.Println("[Env]: get variable:", C.GoString(variable.key))
		return false
	default:
		//fmt.Println("[Env]: command not implemented", cmd)
		return false
	}
	return true
}
