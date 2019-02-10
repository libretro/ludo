// Package savestates takes care of serializing and unserializing the game RAM
// to the host filesystem.
package savestates

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

func name() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	date := time.Now().Format("2006-01-02-15-04-05")
	return name + "@" + date + ".state"
}

// Save the current state to the filesystem
func Save(vid *video.Video) error {
	vid.TakeScreenshot()

	s := state.Global.Core.SerializeSize()
	bytes, err := state.Global.Core.Serialize(s)
	if err != nil {
		return err
	}
	path := filepath.Join(settings.Current.SavestatesDirectory, name())
	err = os.MkdirAll(settings.Current.SavestatesDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// Load the state from the filesystem
func Load(path string) error {
	s := state.Global.Core.SerializeSize()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = state.Global.Core.Unserialize(bytes, s)
	return err
}
