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

// Save the current state to the filesystem. name is the name of the
// savestate file to save to, without extension.
func Save(name string) error {
	s := state.Core.SerializeSize()
	bytes, err := state.Core.Serialize(s)
	if err != nil {
		return err
	}
	path := filepath.Join(settings.Current.SavestatesDirectory, name+".state")
	err = os.MkdirAll(settings.Current.SavestatesDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// Load the state from the filesystem
func Load(path string) error {
	s := state.Core.SerializeSize()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = state.Core.Unserialize(bytes, s)
	return err
}
