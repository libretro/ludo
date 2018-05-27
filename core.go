package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kivutar/go-playthemall/libretro"
)

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

func coreGetGameInfo(filename string, blockExtract bool) (libretro.GameInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return libretro.GameInfo{}, err
	}

	fi, err := file.Stat()
	if err != nil {
		return libretro.GameInfo{}, err
	}

	var realfilename string
	var size int64
	if filepath.Ext(filename) == ".zip" && !blockExtract {
		r, err := zip.OpenReader(filename)
		if err != nil {
			return libretro.GameInfo{}, err
		}
		defer r.Close()

		cf := r.File[0]
		rc, err := cf.Open()
		if err != nil {
			return libretro.GameInfo{}, err
		}
		defer rc.Close()

		f2, err := os.Create("/tmp/" + cf.Name)
		if err != nil {
			return libretro.GameInfo{}, err
		}
		defer f2.Close()
		io.CopyN(f2, rc, int64(cf.UncompressedSize))
		realfilename = "/tmp/" + cf.Name
		size = int64(cf.UncompressedSize)
	} else {
		realfilename = filename
		size = fi.Size()
	}

	gi := libretro.GameInfo{
		Path: realfilename,
		Size: size,
	}

	return gi, nil
}

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
