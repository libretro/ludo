package libretro

import (
	"syscall"
	"unsafe"
)

// Core is an instance of a dynalically loaded libretro core
type Core struct {
	handle *syscall.DLL

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
	proc := core.handle.MustFindProc(name)
	return unsafe.Pointer(proc.Addr())
}

// DlOpen opens a dynamic library
func DlOpen(path string) (*syscall.DLL, error) {
	return syscall.LoadDLL(path)
}
