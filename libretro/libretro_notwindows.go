// +build !windows

package libretro

import (
	"errors"
	"unsafe"
)

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
*/
import "C"

// Core is an instance of a dynalically loaded libretro core
type Core struct {
	handle unsafe.Pointer

	symRetroInit                unsafe.Pointer
	symRetroDeinit              unsafe.Pointer
	symRetroAPIVersion          unsafe.Pointer
	symRetroGetSystemInfo       unsafe.Pointer
	symRetroGetSystemAVInfo     unsafe.Pointer
	symRetroSetEnvironment      unsafe.Pointer
	symRetroSetVideoRefresh     unsafe.Pointer
	symRetroSetInputPoll        unsafe.Pointer
	symRetroSetInputState       unsafe.Pointer
	symRetroSetAudioSample      unsafe.Pointer
	symRetroSetAudioSampleBatch unsafe.Pointer
	symRetroRun                 unsafe.Pointer
	symRetroReset               unsafe.Pointer
	symRetroLoadGame            unsafe.Pointer
	symRetroUnloadGame          unsafe.Pointer
	symRetroSerializeSize       unsafe.Pointer
	symRetroSerialize           unsafe.Pointer
	symRetroUnserialize         unsafe.Pointer

	videoRefresh videoRefreshFunc
}

// DlSym loads a symbol from a dynamic library
func (core *Core) DlSym(name string) unsafe.Pointer {
	return C.dlsym(core.handle, C.CString(name))
}

// DlOpen opens a dynamic library
func DlOpen(path string) (unsafe.Pointer, error) {
	h := C.dlopen(C.CString(path), C.RTLD_NOW)
	if h == nil {
		return h, errors.New("dlopen failed")
	}
	return h, nil
}
