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
	return time.Now().UnixNano() / 1000
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

func environmentSetVariables(data unsafe.Pointer) bool {
	variables := libretro.GetVariables(data)

	pass := []options.VariableInterface{}
	for _, va := range variables {
		va := va
		pass = append(pass, &va)
	}

	var err error
	Options, err = options.New(pass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func environmentSetCoreOptions(data unsafe.Pointer) bool {
	optionDefinitions := libretro.GetCoreOptionDefinitions(data)

	pass := []options.VariableInterface{}
	for _, cod := range optionDefinitions {
		cod := cod
		pass = append(pass, &cod)
	}

	var err error
	Options, err = options.New(pass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func environmentSetCoreOptionsIntl(data unsafe.Pointer) bool {
	optionDefinitions := libretro.GetCoreOptionsIntl(data)

	pass := []options.VariableInterface{}
	for _, cod := range optionDefinitions {
		cod := cod
		pass = append(pass, &cod)
	}

	var err error
	Options, err = options.New(pass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func environment(cmd uint32, data unsafe.Pointer) bool {
	switch cmd {
	case libretro.EnvironmentSetRotation:
		return vid.SetRotation(*(*uint)(data))
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
	case libretro.EnvironmentGetCoreOptionsVersion:
		libretro.SetUint(data, 1)
	case libretro.EnvironmentSetCoreOptions:
		return environmentSetCoreOptions(data)
	case libretro.EnvironmentSetCoreOptionsIntl:
		return environmentSetCoreOptionsIntl(data)
	case libretro.EnvironmentGetVariable:
		return environmentGetVariable(data)
	case libretro.EnvironmentSetVariables:
		return environmentSetVariables(data)
	case libretro.EnvironmentGetVariableUpdate:
		libretro.SetBool(data, Options.Updated)
		Options.Updated = false
	case libretro.EnvironmentSetGeometry:
		vid.Geom = libretro.GetGeometry(data)
	case libretro.EnvironmentSetSystemAVInfo:
		avi := libretro.GetSystemAVInfo(data)
		vid.Geom = avi.Geometry
	case libretro.EnvironmentGetFastforwarding:
		libretro.SetBool(data, state.Global.FastForward)
	case libretro.EnvironmentGetLanguage:
		libretro.SetUint(data, 0)
	default:
		//log.Println("[Env]: Not implemented:", cmd)
		return false
	}
	return true
}
