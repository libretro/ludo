package libretro

import (
	"syscall"
	"unsafe"
)

// DlHandle is a handle to a dynamic library
type DlHandle = *syscall.DLL

// DlSym loads a symbol from a dynamic library
func DlSym(handle DlHandle, name string) unsafe.Pointer {
	proc := handle.MustFindProc(name)
	return unsafe.Pointer(proc.Addr())
}

// DlOpen opens a dynamic library
func DlOpen(path string) (DlHandle, error) {
	return syscall.LoadDLL(path)
}

// DlClose closes a dynamic library
func DlClose(handle DlHandle) error {
	return handle.Release()
}
