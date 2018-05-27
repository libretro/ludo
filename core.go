package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kivutar/go-playthemall/libretro"
)

// coreLoad loads a libretro core
func coreLoad(sofile string) {
	if g.coreRunning {
		g.core.UnloadGame()
		g.core.Deinit()
		g.gamePath = ""
	}

	g.core, _ = libretro.Load(sofile)
	g.core.SetEnvironment(environment)
	g.core.SetVideoRefresh(videoRefresh)
	g.core.SetInputPoll(inputPoll)
	g.core.SetInputState(inputState)
	g.core.SetAudioSample(audioSample)
	g.core.SetAudioSampleBatch(audioSampleBatch)
	g.core.Init()

	// Append the library name to the window title.
	si := g.core.GetSystemInfo()
	if len(si.LibraryName) > 0 {
		if window != nil {
			window.SetTitle("Play Them All - " + si.LibraryName)
		}
		fmt.Println("[Libretro]: Name:", si.LibraryName)
		fmt.Println("[Libretro]: Version:", si.LibraryVersion)
		fmt.Println("[Libretro]: Valid extensions:", si.ValidExtensions)
		fmt.Println("[Libretro]: Need fullpath:", si.NeedFullpath)
		fmt.Println("[Libretro]: Block extract:", si.BlockExtract)
	}

	notify("Core loaded: "+si.LibraryName, 240)
	fmt.Println("[Libretro]: API version:", g.core.APIVersion())
}

// coreUnzipGame unzips a rom to tmpdir and returns the path and size of the extracted rom
func coreUnzipGame(filename string) (string, int64, error) {
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

// coreGetGameInfo opens a rom and return the libretro.GameInfo needed to launch it
func coreGetGameInfo(filename string, blockExtract bool) (libretro.GameInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return libretro.GameInfo{}, err
	}

	fi, err := file.Stat()
	if err != nil {
		return libretro.GameInfo{}, err
	}

	if filepath.Ext(filename) == ".zip" && !blockExtract {
		path, size, err := coreUnzipGame(filename)
		if err != nil {
			return libretro.GameInfo{}, err
		}
		return libretro.GameInfo{Path: path, Size: size}, nil
	}
	return libretro.GameInfo{Path: filename, Size: fi.Size()}, nil
}

// coreLoadGame loads a game. A core has to be loaded first.
func coreLoadGame(filename string) {
	si := g.core.GetSystemInfo()

	gi, err := coreGetGameInfo(filename, si.BlockExtract)
	if err != nil {
		notify(err.Error(), 240)
		fmt.Println(err)
		return
	}

	if !si.NeedFullpath {
		bytes, err := slurp(gi.Path)
		if err != nil {
			notify(err.Error(), 240)
			fmt.Println(err)
			return
		}
		gi.SetData(bytes)
	}

	ok := g.core.LoadGame(gi)
	if !ok {
		notify("The core failed to load the content.", 240)
		fmt.Println("[Libretro]: The core failed to load the content.")
		g.coreRunning = false
		return
	}

	avi := g.core.GetSystemAVInfo()

	// Create the video window
	videoConfigure(avi.Geometry, settings.VideoFullscreen)

	// Append the library name to the window title.
	if len(si.LibraryName) > 0 {
		window.SetTitle("Play Them All - " + si.LibraryName)
	}

	inputInit()
	audioInit(int32(avi.Timing.SampleRate))
	if g.audioCb.SetState != nil {
		g.audioCb.SetState(true)
	}

	g.coreRunning = true
	g.menuActive = false
	g.gamePath = filename
	menuInit()
	notify("Game loaded: "+filename, 240)
	fmt.Println("[Libretro]: Game loaded: " + filename)
}
