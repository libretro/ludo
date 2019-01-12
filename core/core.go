// Package core takes care of instanciating the libretro core, setting the
// input, audio, video, environment callbacks needed to play the games.
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
	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/savefiles"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

// MenuInterface allows passing a *menu.Menu to the core package while avoiding
// cyclic dependencies.
type MenuInterface interface {
	ContextReset()
	UpdateOptions(*options.Options)
}

var vid *video.Video
var opts *options.Options
var menu MenuInterface

// Init is there mainly for dependency injection.
// Call Init before calling other functions of this package.
func Init(v *video.Video, m MenuInterface) {
	vid = v
	menu = m
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			if state.Global.CoreRunning && !state.Global.MenuActive {
				savefiles.SaveSRAM()
			}
		}
	}()
}

// Load loads a libretro core
func Load(sofile string) error {
	state.Global.CorePath = sofile

	// In case the a core is already loaded, we need to close it properly
	// before loading the new core
	Unload()

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

	menu.UpdateOptions(opts)

	log.Println("[Core]: Core loaded: " + si.LibraryName)
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
func LoadGame(filename string) error {

	// If we're loading a new game on the same core, save the RAM of the previous
	// game before closing it.
	if state.Global.GamePath != filename {
		UnloadGame()
	}

	si := state.Global.Core.GetSystemInfo()

	gi, err := getGameInfo(filename, si.BlockExtract)
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

	input.Init(vid, menu)
	audio.Init(int32(avi.Timing.SampleRate))
	if state.Global.AudioCb.SetState != nil {
		state.Global.AudioCb.SetState(true)
	}

	state.Global.CoreRunning = true
	state.Global.MenuActive = false
	state.Global.GamePath = filename

	log.Println("[Core]: Game loaded: " + filename)
	savefiles.LoadSRAM()

	return nil
}

// Unload unloads a libretro core
func Unload() {
	if state.Global.CoreRunning {
		UnloadGame()
		state.Global.Core.Deinit()
	}
}

// UnloadGame unloads a game.
func UnloadGame() {
	savefiles.SaveSRAM()
	state.Global.Core.UnloadGame()
	state.Global.GamePath = ""
	state.Global.CoreRunning = false
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
