package main

import (
	"fmt"
	"libretro"
	"os/user"
	"unsafe"
)

func environment(cmd uint32, data unsafe.Pointer) bool {
	switch cmd {
	case libretro.EnvironmentGetUsername:
		currentUser, err := user.Current()
		if err != nil {
			libretro.SetString(data, "")
		} else {
			libretro.SetString(data, currentUser.Username)
		}
		break
	case libretro.EnvironmentGetLogInterface:
		core.BindLogCallback(data, nanoLog)
		break
	case libretro.EnvironmentGetCanDupe:
		libretro.SetBool(data, true)
		break
	case libretro.EnvironmentSetPixelFormat:
		format := libretro.GetPixelFormat(data)
		if format > libretro.PixelFormatRGB565 {
			return false
		}
		return videoSetPixelFormat(format)
	case libretro.EnvironmentGetSystemDirectory:
	case libretro.EnvironmentGetSaveDirectory:
		libretro.SetString(data, ".")
		return true
	case libretro.EnvironmentShutdown:
		window.SetShouldClose(true)
		return true
	case libretro.EnvironmentGetVariable:
		variable := libretro.GetVariable(data)
		fmt.Println("[Env]: get variable:", variable.Key)
		return false
	default:
		//fmt.Println("[Env]: command not implemented", cmd)
		return false
	}
	return true
}
