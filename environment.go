package main

import (
	"fmt"
	"libretro"
	"os/user"
	"time"
	"unsafe"
)

func getTimeUsec() int64 {
	//fmt.Printf("Seconds since epoch %d", time.Now().Unix())
	return time.Now().UnixNano()
}

var myftc libretro.FrameTimeCallback
var myauc libretro.AudioCallback

func environment(cmd uint32, data unsafe.Pointer) bool {
	switch cmd {
	case libretro.EnvironmentGetUsername:
		currentUser, err := user.Current()
		if err != nil {
			libretro.SetString(data, "")
		} else {
			libretro.SetString(data, currentUser.Username)
		}
	case libretro.EnvironmentGetLogInterface:
		g.core.BindLogCallback(data, nanoLog)
	case libretro.EnvironmentGetPerfInterface:
		g.core.BindPerfCallback(data, getTimeUsec)
	case libretro.EnvironmentSetFrameTimeCallback:
		myftc = libretro.SetFrameTimeCallback(data)
	case libretro.EnvironmentSetAudioCallback:
		myauc = libretro.SetAudioCallback(data)
	case libretro.EnvironmentGetCanDupe:
		libretro.SetBool(data, true)
	case libretro.EnvironmentSetPixelFormat:
		format := libretro.GetPixelFormat(data)
		if format > libretro.PixelFormatRGB565 {
			return false
		}
		return videoSetPixelFormat(format)
	case libretro.EnvironmentGetSystemDirectory:
		fallthrough
	case libretro.EnvironmentGetSaveDirectory:
		libretro.SetString(data, ".")
	case libretro.EnvironmentShutdown:
		window.SetShouldClose(true)
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
