// Package savefiles takes care of saving the game SRAM to the filesystem
package savefiles

import (
	"C"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

var mutex sync.Mutex

func name() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".srm"
}

// path returns the path of the SRAM file for the current core
func path() string {
	return filepath.Join(settings.Current.SavefilesDirectory, name())
}

// SaveSRAM saves the game SRAM to the filesystem
func SaveSRAM() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !state.Global.CoreRunning {
		return errors.New("core not running")
	}

	len := state.Global.Core.GetMemorySize(libretro.MemorySaveRAM)
	ptr := state.Global.Core.GetMemoryData(libretro.MemorySaveRAM)
	if ptr == nil || len == 0 {
		return errors.New("unable to get SRAM address")
	}

	// convert the C array to a go slice
	bytes := C.GoBytes(ptr, C.int(len))
	err := os.MkdirAll(settings.Current.SavefilesDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	fd, err := os.Create(path())
	if err != nil {
		return err
	}
	defer fd.Close()
	fd.Write(bytes)

	return nil
}

// LoadSRAM saves the game SRAM to the filesystem
func LoadSRAM() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !state.Global.CoreRunning {
		return errors.New("core not running")
	}

	fd, err := os.Open(path())
	if err != nil {
		return err
	}
	defer fd.Close()

	len := state.Global.Core.GetMemorySize(libretro.MemorySaveRAM)
	ptr := state.Global.Core.GetMemoryData(libretro.MemorySaveRAM)
	if ptr == nil || len == 0 {
		return errors.New("unable to get SRAM address")
	}

	// this *[1 << 30]byte points to the same memory as ptr, allowing to
	// overwrite this memory
	destination := (*[1 << 30]byte)(unsafe.Pointer(ptr))[:len:len]
	source, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	copy(destination, source)

	return nil
}
