//go:build gltest && linux

package video

import (
	"bytes"
	"runtime"
	"testing"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/libretro/ludo/libretro"
	"golang.org/x/sys/unix"
)

func TestSoftwareFrameUploadBounds(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := glfw.Init(); err != nil {
		t.Fatal(err)
	}
	defer glfw.Terminate()
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.False)
	w, err := glfw.CreateWindow(32, 32, "Ludo upload test", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Destroy()
	w.MakeContextCurrent()
	if err := initGL(); err != nil {
		t.Fatal(err)
	}

	v := &Video{Geom: libretro.GameGeometry{MaxWidth: 64, MaxHeight: 64}}
	v.SetPixelFormat(libretro.PixelFormatXRGB8888)
	v.resetGameTexture(1, 1)
	defer gl.DeleteTextures(1, &v.texID)
	v.program, err = newProgram(vertexShader, defaultFragmentShader)
	if err != nil {
		t.Fatal(err)
	}
	defer gl.DeleteProgram(v.program)

	const width, height, pitch = 32, 32, 32 * 4
	page := unix.Getpagesize()
	memory, err := unix.Mmap(-1, 0, page*2, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Munmap(memory)
	if err := unix.Mprotect(memory[page:], unix.PROT_NONE); err != nil {
		t.Fatal(err)
	}
	// The core supplies only the current frame, ending at a protected page.
	frame := memory[page-pitch*height : page]
	for i := range frame {
		frame[i] = byte(i)
	}
	v.Refresh(unsafe.Pointer(&frame[0]), width, height, pitch)
	v.uploadTexture()

	var texWidth, texHeight int32
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &texWidth)
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_HEIGHT, &texHeight)
	if texWidth != width || texHeight != height {
		t.Fatalf("texture size = %dx%d, want %dx%d", texWidth, texHeight, width, height)
	}
	pixels := make([]byte, len(frame))
	gl.GetTexImage(gl.TEXTURE_2D, 0, v.pixType, v.pixFmt, gl.Ptr(pixels))
	if !bytes.Equal(pixels, frame) {
		t.Error("uploaded pixels differ from the supplied frame")
	}
	if err := gl.GetError(); err != gl.NO_ERROR {
		t.Errorf("GL error: %#x", err)
	}
}
