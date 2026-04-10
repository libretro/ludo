// Package core takes care of instanciating the libretro core, setting the
// input, audio, video, environment callbacks needed to play the games.
// It also deals with core options and persisting SRAM periodically.
package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/patch"
	"github.com/libretro/ludo/savefiles"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"

	"github.com/mholt/archives"
)

var vid *video.Video

// Options holds the settings for the current core
var Options *options.Options

// Init is there mainly for dependency injection.
// Call Init before calling other functions of this package.
func Init(v *video.Video) {
	vid = v
}

// Load loads a libretro core
func Load(sofile string) error {
	// In case a core is already loaded, we need to close it properly
	// before loading the new core
	Unload()

	// This must be set before the environment callback is called
	state.CorePath = sofile

	var err error
	state.Core, err = libretro.Load(sofile)
	if err != nil {
		return err
	}
	state.Core.SetEnvironment(environment)
	state.Core.Init()
	state.Core.SetVideoRefresh(vid.Refresh)
	state.Core.SetInputPoll(func() {})
	state.Core.SetInputState(input.State)
	state.Core.SetAudioSample(audio.Sample)
	state.Core.SetAudioSampleBatch(audio.SampleBatch)

	// Append the library name to the window title.
	si := state.Core.GetSystemInfo()
	if len(si.LibraryName) > 0 {
		vid.SetTitle("Ludo - " + si.LibraryName)
		if state.Verbose {
			log.Println("[Core]: Name:", si.LibraryName)
			log.Println("[Core]: Version:", si.LibraryVersion)
			log.Println("[Core]: Valid extensions:", si.ValidExtensions)
			log.Println("[Core]: Need fullpath:", si.NeedFullpath)
			log.Println("[Core]: Block extract:", si.BlockExtract)
		}
	}

	return nil
}

// unarchiveGame unarchives a rom to tmpdir and returns the path and size of the extracted ROM.
// In case the archive contains more than one file, they are all extracted and the
// first one or a better match (cue for CDrom) is passed to the libretro core.
func unarchiveGame(filename string) (string, int64, error) {
	path := ""
	size := int64(0)
	dst := os.TempDir()
	ctx := context.Background()

	file, err := os.Open(filename)
	if err != nil {
		return path, size, err
	}
	defer file.Close()

	format, stream, err := archives.Identify(ctx, filepath.Base(filename), file)
	if err != nil {
		return path, size, err
	}

	extPriority := 0
	extPrefered := map[string]int{
		".cue": 1,
		".m3u": 2,
		".pbp": 3,
	}

	selectGame := func(fname string, fsize int64) {
		ext := filepath.Ext(fname)
		if size == 0 {
			path = filepath.Join(dst, fname)
			size = fsize
			log.Println("first file in archive:", path, size)
		}
		priority, ok := extPrefered[strings.ToLower(ext)]
		if ok && priority > extPriority {
			extPriority = priority
			path = filepath.Join(dst, fname)
			size = fsize
			log.Println("find a better file in archive:", path, size)
		}
	}

	if ex, ok := format.(archives.Extractor); ok {
		err = ex.Extract(ctx, stream, func(ctx context.Context, f archives.FileInfo) error {
			fname, err := archiveOutputPath(dst, f.NameInArchive)
			if err != nil {
				return err
			}
			if f.IsDir() {
				return os.MkdirAll(fname, os.ModePerm)
			}
			if err := os.MkdirAll(filepath.Dir(fname), os.ModePerm); err != nil {
				return err
			}

			in, err := f.Open()
			if err != nil {
				return err
			}
			defer in.Close()

			mode := f.Mode().Perm()
			if mode == 0 {
				mode = 0644
			}
			out, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return err
			}

			_, copyErr := io.Copy(out, in)
			closeErr := out.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}

			selectGame(strings.TrimPrefix(fname, dst+string(filepath.Separator)), f.Size())
			return nil
		})
		return path, size, err
	}

	decomp, ok := format.(archives.Decompressor)
	if !ok {
		return path, size, fmt.Errorf("unsupported archive format: %T", format)
	}

	rc, err := decomp.OpenReader(stream)
	if err != nil {
		return path, size, err
	}
	defer rc.Close()

	fname := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	if fname == "" {
		return path, size, errors.New("could not determine extracted filename")
	}
	path = filepath.Join(dst, fname)
	out, err := os.Create(path)
	if err != nil {
		return "", 0, err
	}
	size, err = io.Copy(out, rc)
	closeErr := out.Close()
	if err != nil {
		return "", 0, err
	}
	if closeErr != nil {
		return "", 0, closeErr
	}

	return path, size, nil
}

func archiveOutputPath(dst, name string) (string, error) {
	clean := path.Clean(strings.TrimPrefix(name, "/"))
	if clean == "." || clean == "" {
		return "", errors.New("invalid archive path")
	}
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("unsafe archive path: %s", name)
	}
	fname := filepath.Join(dst, filepath.FromSlash(clean))
	rel, err := filepath.Rel(dst, fname)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe archive path: %s", name)
	}
	return fname, nil
}

// LoadGame loads a game. A core has to be loaded first.
func LoadGame(gamePath string) error {
	if _, err := os.Stat(gamePath); os.IsNotExist(err) {
		return err
	}

	// If we're loading a new game on the same core, save the RAM of the previous
	// game before closing it.
	if state.GamePath != gamePath {
		UnloadGame()
	}

	if state.Core == nil {
		return errors.New("no core loaded")
	}

	si := state.Core.GetSystemInfo()

	gi, err := getGameInfo(gamePath, si.BlockExtract)
	if err != nil {
		return err
	}

	if !si.NeedFullpath {
		bytes, err := os.ReadFile(gi.Path)
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

	ok := state.Core.LoadGame(*gi)
	if !ok {
		state.CoreRunning = false
		return errors.New("failed to load the game")
	}

	avi := state.Core.GetSystemAVInfo()

	vid.Geom = avi.Geometry

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		vid.SetTitle("Ludo - " + si.LibraryName)
	}

	input.Init(vid)
	audio.Reconfigure(int32(avi.Timing.SampleRate))
	if state.Core.AudioCallback != nil {
		state.Core.AudioCallback.SetState(true)
	}

	state.CoreRunning = true
	state.FastForward = false
	state.GamePath = gamePath

	state.Core.SetControllerPortDevice(0, libretro.DeviceJoypad)
	state.Core.SetControllerPortDevice(1, libretro.DeviceJoypad)
	state.Core.SetControllerPortDevice(2, libretro.DeviceJoypad)
	state.Core.SetControllerPortDevice(3, libretro.DeviceJoypad)
	state.Core.SetControllerPortDevice(4, libretro.DeviceJoypad)

	log.Println("[Core]: Game loaded: " + gamePath)
	savefiles.LoadSRAM()

	return nil
}

// Unload unloads a libretro core
func Unload() {
	if state.Core != nil {
		UnloadGame()
		state.Core.Deinit()
		state.CorePath = ""
		state.Core = nil
		Options = nil
	}
}

// UnloadGame unloads a game.
func UnloadGame() {
	if state.CoreRunning {
		savefiles.SaveSRAM()
		state.Core.UnloadGame()
		state.GamePath = ""
		state.CoreRunning = false
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

	if !blockExtract {
		switch filepath.Ext(filename) {
		case ".zip", ".zst", ".rar", ".tar":
			path, size, err := unarchiveGame(filename)
			if err != nil {
				return nil, err
			}
			return &libretro.GameInfo{Path: path, Size: size}, nil
		}
	}

	return &libretro.GameInfo{Path: filename, Size: fi.Size()}, nil
}
