// Package savestates takes care of serializing and unserializing the game RAM
// to the host filesystem.
package savestates

import (
	"os"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/settings"
)

// Save the current state to the filesystem. name is the name of the
// savestate file to save to, without extension.
func Save(name string) error {
	bytes, err := core.Serialize()
	if err != nil {
		return err
	}
	path := filepath.Join(settings.Current.SavestatesDirectory, name+".state")
	err = os.MkdirAll(settings.Current.SavestatesDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// Load the state from the filesystem
func Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return core.Unserialize(bytes)
}
