package video

import (
	"testing"

	"github.com/libretro/ludo/libretro"
)

func TestDuplicateFramePreservesDimensions(t *testing.T) {
	v := &Video{}
	v.Refresh(libretro.HWFrameBufferValid, 320, 240, 0)
	v.Refresh(nil, 0, 0, 0)
	if v.width != 320 || v.height != 240 || v.data != libretro.HWFrameBufferValid || !v.needUpload {
		t.Fatal("NULL refresh discarded the pending frame")
	}
	v.needUpload = false
	v.Refresh(nil, 640, 480, 0)
	if v.width != 320 || v.height != 240 || v.needUpload {
		t.Fatal("NULL refresh changed the previously presented frame")
	}
}

func TestContextDestroyWithoutCore(t *testing.T) {
	(&Video{}).RunContextDestroy()
}
