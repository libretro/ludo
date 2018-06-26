package main

import (
	"log"
	"os/user"
	"time"
	"unsafe"

	"github.com/kivutar/go-playthemall/libretro"
)

var logLevels = map[uint32]string{
	libretro.LogLevelDebug: "DEBUG",
	libretro.LogLevelInfo:  "INFO",
	libretro.LogLevelWarn:  "WARN",
	libretro.LogLevelError: "ERROR",
	libretro.LogLevelDummy: "DUMMY",
}

func logCallback(level uint32, str string) {
	log.Printf("[%s]: %s", logLevels[level], str)
}

func getTimeUsec() int64 {
	//fmt.Printf("Seconds since epoch %d", time.Now().Unix())
	return time.Now().UnixNano()
}

var options struct {
	Updated bool
	Vars    []libretro.Variable
}

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
		g.core.BindLogCallback(data, logCallback)
	case libretro.EnvironmentGetPerfInterface:
		g.core.BindPerfCallback(data, getTimeUsec)
	case libretro.EnvironmentSetFrameTimeCallback:
		g.frameTimeCb = libretro.SetFrameTimeCallback(data)
	case libretro.EnvironmentSetAudioCallback:
		g.audioCb = libretro.SetAudioCallback(data)
	case libretro.EnvironmentGetCanDupe:
		libretro.SetBool(data, true)
	case libretro.EnvironmentSetPixelFormat:
		format := libretro.GetPixelFormat(data)
		if format > libretro.PixelFormatRGB565 {
			return false
		}
		return videoSetPixelFormat(format)
	case libretro.EnvironmentGetSystemDirectory:
		usr, _ := user.Current()
		libretro.SetString(data, usr.HomeDir+"/.playthemall/system/")
	case libretro.EnvironmentGetSaveDirectory:
		usr, _ := user.Current()
		libretro.SetString(data, usr.HomeDir+"/.playthemall/savefiles/")
	case libretro.EnvironmentShutdown:
		window.SetShouldClose(true)
	case libretro.EnvironmentGetVariable:
		variable := libretro.GetVariable(data)
		if variable.Key() == "snes9x_layer_1" {
			variable.SetValue("disabled")
			return true
		}
		return false
	case libretro.EnvironmentSetVariables:
		options.Vars = libretro.GetVariables(data)
		options.Updated = true
		return true
	case libretro.EnvironmentGetVariableUpdate:
		libretro.SetBool(data, options.Updated)
		options.Updated = false
		return true
	default:
		//log.Println("[Env]: Not implemented:", cmd)
		return false
	}
	return true
}
