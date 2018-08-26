package core

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/libretro/go-playthemall/audio"
	"github.com/libretro/go-playthemall/input"
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/notifications"
	"github.com/libretro/go-playthemall/options"
	"github.com/libretro/go-playthemall/state"
	"github.com/libretro/go-playthemall/utils"
	"github.com/libretro/go-playthemall/video"
)

type ContextReseter interface {
	ContextReset()
}

var vid *video.Video
var opts *options.Options
var menu ContextReseter

// Init is there mainly for dependency injection.
// Call Init before calling other functions of this package.
func Init(v *video.Video, o *options.Options, m ContextReseter) {
	vid = v
	opts = o
	menu = m
}

// Load loads a libretro core
func Load(sofile string) {
	state.Global.CorePath = sofile
	if state.Global.CoreRunning {
		state.Global.Core.UnloadGame()
		state.Global.Core.Deinit()
		state.Global.GamePath = ""
		state.Global.CoreRunning = false
	}

	state.Global.Core, _ = libretro.Load(sofile)
	state.Global.Core.SetEnvironment(environment)
	state.Global.Core.SetVideoRefresh(vid.Refresh)
	state.Global.Core.SetInputPoll(input.Poll)
	state.Global.Core.SetInputState(input.State)
	state.Global.Core.SetAudioSample(audio.Sample)
	state.Global.Core.SetAudioSampleBatch(audio.SampleBatch)
	state.Global.Core.Init()

	// Append the library name to the window title.
	si := state.Global.Core.GetSystemInfo()
	if len(si.LibraryName) > 0 {
		if vid.Window != nil {
			vid.Window.SetTitle("Play Them All - " + si.LibraryName)
		}
		if state.Global.Verbose {
			log.Println("[Core]: Name:", si.LibraryName)
			log.Println("[Core]: Version:", si.LibraryVersion)
			log.Println("[Core]: Valid extensions:", si.ValidExtensions)
			log.Println("[Core]: Need fullpath:", si.NeedFullpath)
			log.Println("[Core]: Block extract:", si.BlockExtract)
		}
	}

	notifications.DisplayAndLog("Core", "Core loaded: "+si.LibraryName)
}

// unzipGame unzips a rom to tmpdir and returns the path and size of the extracted rom
func unzipGame(filename string) (string, int64, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", 0, err
	}
	defer r.Close()

	cf := r.File[0]
	size := int64(cf.UncompressedSize)
	rc, err := cf.Open()
	if err != nil {
		return "", 0, err
	}
	defer rc.Close()

	path := os.TempDir() + "/" + cf.Name

	f2, err := os.Create(path)
	if err != nil {
		return "", 0, err
	}
	defer f2.Close()
	_, err = io.CopyN(f2, rc, size)
	if err != nil {
		return "", 0, err
	}

	return path, size, nil
}

// LoadGame loads a game. A core has to be loaded first.
func LoadGame(filename string) {
	si := state.Global.Core.GetSystemInfo()

	gi, err := getGameInfo(filename, si.BlockExtract)
	if err != nil {
		notifications.DisplayAndLog("Core", err.Error())
		return
	}

	if !si.NeedFullpath {
		bytes, err := utils.Slurp(gi.Path)
		if err != nil {
			notifications.DisplayAndLog("Core", err.Error())
			return
		}
		gi.SetData(bytes)
	}

	ok := state.Global.Core.LoadGame(gi)
	if !ok {
		notifications.DisplayAndLog("Core", "Failed to load the content.")
		state.Global.CoreRunning = false
		return
	}

	avi := state.Global.Core.GetSystemAVInfo()

	vid.Geom = avi.Geometry

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		vid.Window.SetTitle("Play Them All - " + si.LibraryName)
	}

	input.Init(vid, menu)
	audio.Init(int32(avi.Timing.SampleRate))
	if state.Global.AudioCb.SetState != nil {
		state.Global.AudioCb.SetState(true)
	}

	state.Global.CoreRunning = true
	state.Global.MenuActive = false
	state.Global.GamePath = filename

	notifications.DisplayAndLog("Core", "Game loaded: "+filename)
}

// getGameInfo opens a rom and return the libretro.GameInfo needed to launch it
func getGameInfo(filename string, blockExtract bool) (libretro.GameInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return libretro.GameInfo{}, err
	}

	fi, err := file.Stat()
	if err != nil {
		return libretro.GameInfo{}, err
	}

	if filepath.Ext(filename) == ".zip" && !blockExtract {
		path, size, err := unzipGame(filename)
		if err != nil {
			return libretro.GameInfo{}, err
		}
		return libretro.GameInfo{Path: path, Size: size}, nil
	}
	return libretro.GameInfo{Path: filename, Size: fi.Size()}, nil
}
