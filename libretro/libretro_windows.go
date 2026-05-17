package libretro

/*
#include <stdint.h>

static void *symbol_from_addr(uintptr_t addr) {
	return (void *)addr;
}
*/
import "C"
import (
	"syscall"
	"unsafe"
)

// DlHandle is a handle to a dynamic library
type DlHandle = *syscall.DLL

// DlSym loads a symbol from a dynamic library
func DlSym(handle DlHandle, name string) unsafe.Pointer {
	proc := handle.MustFindProc(name)
	return C.symbol_from_addr(C.uintptr_t(proc.Addr()))
}

// DlOpen opens a dynamic library
func DlOpen(path string) (DlHandle, error) {
	return syscall.LoadDLL(path)
}

// DlClose closes a dynamic library
func DlClose(handle DlHandle) error {
	return handle.Release()
}
