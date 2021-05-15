// +build !windows

package libretro

import (
	"errors"
	"unsafe"
)

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

// Core is an instance of a dynalically loaded libretro core
type Core struct {
	handle unsafe.Pointer

	symRetroInit                    unsafe.Pointer
	symRetroDeinit                  unsafe.Pointer
	symRetroAPIVersion              unsafe.Pointer
	symRetroGetSystemInfo           unsafe.Pointer
	symRetroGetSystemAVInfo         unsafe.Pointer
	symRetroSetEnvironment          unsafe.Pointer
	symRetroSetVideoRefresh         unsafe.Pointer
	symRetroSetControllerPortDevice unsafe.Pointer
	symRetroSetInputPoll            unsafe.Pointer
	symRetroSetInputState           unsafe.Pointer
	symRetroSetAudioSample          unsafe.Pointer
	symRetroSetAudioSampleBatch     unsafe.Pointer
	symRetroRun                     unsafe.Pointer
	symRetroReset                   unsafe.Pointer
	symRetroLoadGame                unsafe.Pointer
	symRetroUnloadGame              unsafe.Pointer
	symRetroSerializeSize           unsafe.Pointer
	symRetroSerialize               unsafe.Pointer
	symRetroUnserialize             unsafe.Pointer
	symRetroGetMemorySize           unsafe.Pointer
	symRetroGetMemoryData           unsafe.Pointer

	AudioCallback       *AudioCallback
	FrameTimeCallback   *FrameTimeCallback
	DiskControlCallback *DiskControlCallback
}

// DlSym loads a symbol from a dynamic library
func (core *Core) DlSym(name string) unsafe.Pointer {
	return C.dlsym(core.handle, C.CString(name))
}

// DlOpen opens a dynamic library
func DlOpen(path string) (unsafe.Pointer, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	h := C.dlopen(cpath, C.RTLD_LAZY|C.RTLD_GLOBAL)
	cerr := C.dlerror()
	if h == nil || cerr != nil {
		err := C.GoString(cerr)
		return nil, errors.New(err)
	}
	return h, nil
}
