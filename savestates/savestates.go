// Package savestates takes care of serializing and unserializing the game RAM
// to the host filesystem.
package savestates

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

func name() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".state"
}

// Save the current state to the filesystem
func Save() error {
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
func Load() error {
	s := state.Global.Core.SerializeSize()
	path := filepath.Join(settings.Current.SavestatesDirectory, name())
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = state.Global.Core.Unserialize(bytes, s)
	return err
}
