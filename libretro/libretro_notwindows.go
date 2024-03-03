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

// DlHandle is a handle to a dynamic library
type DlHandle = unsafe.Pointer

// DlSym loads a symbol from a dynamic library
func DlSym(handle DlHandle, name string) unsafe.Pointer {
	return C.dlsym(handle, C.CString(name))
}

// DlOpen opens a dynamic library
func DlOpen(path string) (DlHandle, error) {
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

// DlClose closes a dynamic library
func DlClose(handle DlHandle) error {
	result := C.dlclose(handle)
	if int(result) != 0 {
		cerr := C.dlerror()
		if cerr != nil {
			err := C.GoString(cerr)
			return errors.New(err)
		}
	}
	return nil
}
