// Package core takes care of instanciating the libretro core, setting the
// input, audio, video, environment callbacks needed to play the games.
// It also deals with core options and persisting SRAM periodically.
package core

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/savefiles"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

var vid *video.Video

// Options holds the settings for the current core
var Options *options.Options

// Init is there mainly for dependency injection.
// Call Init before calling other functions of this package.
func Init(v *video.Video) {
	vid = v
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			state.Global.Lock()
			running := state.Global.CoreRunning
			state.Global.Unlock()
			if running && !state.Global.MenuActive {
				savefiles.SaveSRAM()
			}
		}
	}()
}

// Load loads a libretro core
func Load(sofile string) error {
	// In case the a core is already loaded, we need to close it properly
	// before loading the new core
	Unload()

	// This must be set before the environment callback is called
	state.Global.CorePath = sofile

	var err error
	state.Global.Core, err = libretro.Load(sofile)
	if err != nil {
		return err
	}
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
			vid.Window.SetTitle("Ludo - " + si.LibraryName)
		}
		if state.Global.Verbose {
			log.Println("[Core]: Name:", si.LibraryName)
			log.Println("[Core]: Version:", si.LibraryVersion)
			log.Println("[Core]: Valid extensions:", si.ValidExtensions)
			log.Println("[Core]: Need fullpath:", si.NeedFullpath)
			log.Println("[Core]: Block extract:", si.BlockExtract)
		}
	}

	return nil
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
func LoadGame(gamePath string) error {
	// If we're loading a new game on the same core, save the RAM of the previous
	// game before closing it.
	if state.Global.GamePath != gamePath {
		UnloadGame()
	}

	si := state.Global.Core.GetSystemInfo()

	gi, err := getGameInfo(gamePath, si.BlockExtract)
	if err != nil {
		return err
	}

	if !si.NeedFullpath {
		bytes, err := ioutil.ReadFile(gi.Path)
		if err != nil {
			return err
		}
		gi.SetData(bytes)
	}

	ok := state.Global.Core.LoadGame(*gi)
	if !ok {
		state.Global.CoreRunning = false
		return errors.New("failed to load the game")
	}

	avi := state.Global.Core.GetSystemAVInfo()

	vid.Geom = avi.Geometry

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		vid.Window.SetTitle("Ludo - " + si.LibraryName)
	}

	input.Init(vid)
	audio.Reconfigure(int32(avi.Timing.SampleRate))
	if state.Global.Core.AudioCallback != nil {
		state.Global.Core.AudioCallback.SetState(true)
	}

	state.Global.Lock()
	state.Global.CoreRunning = true
	state.Global.FastForward = false
	state.Global.Unlock()

	state.Global.GamePath = gamePath

	state.Global.Core.SetControllerPortDevice(0, libretro.DeviceJoypad)
	state.Global.Core.SetControllerPortDevice(1, libretro.DeviceJoypad)

	log.Println("[Core]: Game loaded: " + gamePath)
	savefiles.LoadSRAM()

	ntf.Display(ntf.Info, "Press P to toggle the menu", ntf.Medium)

	vid.InitFramebuffer(vid.Geom.BaseWidth, vid.Geom.BaseHeight)

	if state.Global.Core.HWRenderCallback != nil {
		state.Global.Core.HWRenderCallback.ContextReset()
	}

	return nil
}

// Unload unloads a libretro core
func Unload() {
	if state.Global.CoreRunning {
		UnloadGame()
		state.Global.Core.Deinit()
		state.Global.CorePath = ""
		state.Global.Core = nil
		Options = nil
	}
}

// UnloadGame unloads a game.
func UnloadGame() {
	if state.Global.CoreRunning {
		savefiles.SaveSRAM()
		state.Global.Core.UnloadGame()
		state.Global.GamePath = ""
		state.Global.CoreRunning = false
		vid.ResetPitch()
	}
}

// getGameInfo opens a rom and return the libretro.GameInfo needed to launch it
func getGameInfo(filename string, blockExtract bool) (*libretro.GameInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if filepath.Ext(filename) == ".zip" && !blockExtract {
		path, size, err := unzipGame(filename)
		if err != nil {
			return nil, err
		}
		return &libretro.GameInfo{Path: path, Size: size}, nil
	}
	return &libretro.GameInfo{Path: filename, Size: fi.Size()}, nil
}
