// Package savefiles takes care of saving the game SRAM to the filesystem
package savefiles

import (
	"os"
	"path/filepath"

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

// SaveSRAM saves the game SRAM to the filesystem
func SaveSRAM() {
	if state.Global.CoreRunning {
		len := state.Global.Core.GetMemorySize(libretro.MemorySaveRAM)
		dat := state.Global.Core.GetMemoryData(libretro.MemorySaveRAM)
		bytes := make([]byte, len)
		copy(bytes, (*(*[]byte)(dat))[:])
		//fmt.Println((*(*[]byte)(dat))[:])
		path := filepath.Join(settings.Current.SavefilesDirectory, name())
		os.MkdirAll(settings.Current.SavefilesDirectory, os.ModePerm)
		srm, _ := os.Create(path)
		srm.Write(bytes)
	}
}
