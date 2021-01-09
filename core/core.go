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

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/patch"
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
	// ticker := time.NewTicker(time.Second * 10)
	// go func() {
	// 	for range ticker.C {
	// 		state.Global.Lock()
	// 		canSave := state.Global.CoreRunning && !state.Global.MenuActive
	// 		state.Global.Unlock()
	// 		if canSave {
	// 			savefiles.SaveSRAM()
	// 		}
	// 	}
	// }()
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
	state.Global.Core.Init()
	state.Global.Core.SetVideoRefresh(vid.Refresh)
	state.Global.Core.SetInputPoll(func() {})
	state.Global.Core.SetInputState(input.State)
	state.Global.Core.SetAudioSample(audio.Sample)
	state.Global.Core.SetAudioSampleBatch(audio.SampleBatch)

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

// unzipGame unzips a rom to tmpdir and returns the path and size of the extracted ROM.
// In case the zip contains more than one file, they are all extracted and the
// first one is passed to the libretro core.
func unzipGame(filename string) (string, int64, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", 0, err
	}
	defer r.Close()

	var mainPath string
	var mainSize int64
	for i, cf := range r.File {
		size := int64(cf.UncompressedSize)
		rc, err := cf.Open()
		if err != nil {
			return "", 0, err
		}

		path := filepath.Join(os.TempDir(), cf.Name)

		if cf.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return "", 0, err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, cf.Mode())
		if err != nil {
			return "", 0, err
		}

		if _, err = io.Copy(outFile, rc); err != nil {
			return "", 0, err
		}
		outFile.Close()
		rc.Close()

		if i == 0 {
			mainPath = path
			mainSize = size
		}
	}

	return mainPath, mainSize, nil
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

		if patched, _ := patch.Try(gamePath, bytes); patched != nil {
			gi.Size = int64(len(*patched))
			gi.SetData(*patched)
		} else {
			gi.SetData(bytes)
		}
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
	state.Global.GamePath = gamePath
	state.Global.Unlock()

	state.Global.Core.SetControllerPortDevice(0, libretro.DeviceJoypad)
	state.Global.Core.SetControllerPortDevice(1, libretro.DeviceJoypad)
	state.Global.Core.SetControllerPortDevice(2, libretro.DeviceJoypad)
	state.Global.Core.SetControllerPortDevice(3, libretro.DeviceJoypad)
	state.Global.Core.SetControllerPortDevice(4, libretro.DeviceJoypad)

	log.Println("[Core]: Game loaded: " + gamePath)
	// savefiles.LoadSRAM()

	return nil
}

// Unload unloads a libretro core
func Unload() {
	if state.Global.Core != nil {
		UnloadGame()
		state.Global.Core.Deinit()
		state.Global.Lock()
		state.Global.CorePath = ""
		state.Global.Core = nil
		Options = nil
		state.Global.Unlock()
	}
}

// UnloadGame unloads a game.
func UnloadGame() {
	if state.Global.CoreRunning {
		//savefiles.SaveSRAM()
		state.Global.Core.UnloadGame()
		state.Global.Lock()
		state.Global.GamePath = ""
		state.Global.CoreRunning = false
		state.Global.Unlock()
		vid.ResetPitch()
		vid.ResetRot()
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
