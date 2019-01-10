// Package savefiles takes care of saving the game SRAM to the filesystem
package savefiles

import (
	"C"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

func name() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".srm"
}

// Path returns the path of the SRAM file for the current core
func Path() string {
	return filepath.Join(settings.Current.SavefilesDirectory, name())
}

// SaveSRAM saves the game SRAM to the filesystem
func SaveSRAM() {
	if state.Global.CoreRunning {
		len := state.Global.Core.GetMemorySize(libretro.MemorySaveRAM)
		ptr := state.Global.Core.GetMemoryData(libretro.MemorySaveRAM)
		if ptr == nil || len == 0 {
			log.Println("[Core]: Unable to get SRAM address")
			return
		}
		// convert the C array to a go slice
		bytes := C.GoBytes(ptr, C.int(len))
		err := os.MkdirAll(settings.Current.SavefilesDirectory, os.ModePerm)
		if err != nil {
			log.Println("[Core]:", err)
			return
		}
		fd, err := os.Create(Path())
		if err != nil {
			log.Println("[Core]:", err)
			return
		}
		fd.Write(bytes)
		log.Println("[Core]: Saved SRAM", Path())
	}
}

// LoadSRAM saves the game SRAM to the filesystem
func LoadSRAM() {
	if state.Global.CoreRunning {
		fd, err := os.Open(Path())
		if err != nil {
			log.Println("[Core]:", err)
			return
		}
		len := state.Global.Core.GetMemorySize(libretro.MemorySaveRAM)
		ptr := state.Global.Core.GetMemoryData(libretro.MemorySaveRAM)
		if ptr == nil || len == 0 {
			log.Println("[Core]: Unable to get SRAM address")
			return
		}
		// this *[1 << 30]byte points to the same memory as ptr, allowing to
		// overwrite this memory
		destination := (*[1 << 30]byte)(unsafe.Pointer(ptr))[:len:len]
		source, err := ioutil.ReadAll(fd)
		if err != nil {
			log.Println("[Core]:", err)
			return
		}
		copy(destination, source)
		log.Println("[Core]: Loaded SRAM", Path())
	}
}
