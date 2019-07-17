package core

import (
	"log"
	"os"
	"os/user"
	"time"
	"unsafe"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
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

func environmentGetVariable(data unsafe.Pointer) bool {
	variable := libretro.GetVariable(data)
	for _, v := range Options.Vars {
		if variable.Key() == v.Key {
			variable.SetValue(v.Choices[v.Choice])
			return true
		}
	}
	return false
}

func environmentSetPixelFormat(data unsafe.Pointer) bool {
	format := libretro.GetPixelFormat(data)
	if format > libretro.PixelFormatRGB565 {
		return false
	}
	return vid.SetPixelFormat(format)
}

func environmentGetUsername(data unsafe.Pointer) bool {
	currentUser, err := user.Current()
	if err != nil {
		libretro.SetString(data, "")
	} else {
		libretro.SetString(data, currentUser.Username)
	}
	return true
}

func environmentGetSystemDirectory(data unsafe.Pointer) bool {
	err := os.MkdirAll(settings.Current.SystemDirectory, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	libretro.SetString(data, settings.Current.SystemDirectory)
	return true
}

func environmentGetSaveDirectory(data unsafe.Pointer) bool {
	err := os.MkdirAll(settings.Current.SavefilesDirectory, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	libretro.SetString(data, settings.Current.SavefilesDirectory)
	return true
}

func environmentSetVariable(data unsafe.Pointer) bool {
	var err error
	Options, err = options.New(libretro.GetVariables(data))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func environment(cmd uint32, data unsafe.Pointer) bool {
	switch cmd {
	case libretro.EnvironmentGetUsername:
		return environmentGetUsername(data)
	case libretro.EnvironmentGetLogInterface:
		state.Global.Core.BindLogCallback(data, logCallback)
	case libretro.EnvironmentGetPerfInterface:
		state.Global.Core.BindPerfCallback(data, getTimeUsec)
	case libretro.EnvironmentSetFrameTimeCallback:
		state.Global.Core.SetFrameTimeCallback(data)
	case libretro.EnvironmentSetAudioCallback:
		state.Global.Core.SetAudioCallback(data)
	case libretro.EnvironmentGetCanDupe:
		libretro.SetBool(data, true)
	case libretro.EnvironmentSetPixelFormat:
		return environmentSetPixelFormat(data)
	case libretro.EnvironmentGetSystemDirectory:
		return environmentGetSystemDirectory(data)
	case libretro.EnvironmentGetSaveDirectory:
		return environmentGetSaveDirectory(data)
	case libretro.EnvironmentShutdown:
		vid.Window.SetShouldClose(true)
	case libretro.EnvironmentGetVariable:
		return environmentGetVariable(data)
	case libretro.EnvironmentSetVariables:
		return environmentSetVariable(data)
	case libretro.EnvironmentGetVariableUpdate:
		libretro.SetBool(data, Options.Updated)
		Options.Updated = false
	case libretro.EnvironmentSetGeometry:
		vid.Geom = libretro.GetGeometry(data)
	case libretro.EnvironmentSetSystemAVInfo:
		avi := libretro.GetSystemAVInfo(data)
		vid.Geom = avi.Geometry
	default:
		//log.Println("[Env]: Not implemented:", cmd)
		return false
	}
	return true
}
