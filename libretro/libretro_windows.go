package libretro

import (
	"syscall"
	"unsafe"
)

// Core is an instance of a dynalically loaded libretro core
type Core struct {
	handle syscall.Handle

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

func (core *Core) DlSym(name string) unsafe.Pointer {
	tmp, _ := syscall.GetProcAddress(core.handle, name)
	return unsafe.Pointer(tmp)
}

func DlOpen(path string) (syscall.Handle, error) {
	return yscall.LoadLibrary(path)
}
