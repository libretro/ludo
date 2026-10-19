//go:build gltest && !darwin

package video

import (
	"runtime"
	"testing"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/state"
)

// Run with go test -tags 'wayland gltest' ./video on a graphical session.
func TestHardwareContextIsolation(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := glfw.Init(); err != nil {
		t.Fatal(err)
	}
	defer glfw.Terminate()
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	w, err := glfw.CreateWindow(128, 96, "Ludo GL test", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Destroy()
	w.MakeContextCurrent()
	if err := initGL(); err != nil {
		t.Fatal(err)
	}
	t.Log("GL:", gl.GoStr(gl.GetString(gl.VERSION)))

	oldCore, oldRunning := state.Core, state.CoreRunning
	defer func() { state.Core, state.CoreRunning = oldCore, oldRunning }()
	hw := &libretro.HWRenderCallback{HWContextType: libretro.HWContextOpenGLCore, VersionMajor: 3, VersionMinor: 3}
	state.Core = &libretro.Core{HWRenderCallback: hw}
	state.CoreRunning = true
	v := &Video{Window: w, identityMat: mgl32.Ident4(), Geom: libretro.GameGeometry{MaxWidth: 64, MaxHeight: 64}}
	v.SetHWRenderContext(hw)
	v.ensurePixelFormatDefaults()
	v.resetGameTexture(64, 64)
	hw.ContextReset = func() {
		if glfw.GetCurrentContext() != v.HWWindow {
			t.Fatal("context_reset ran outside the HW context")
		}
	}
	if err := v.ResetHWContext(); err != nil {
		t.Fatal(err)
	}
	defer v.DestroySharedHWContext()
	if glfw.GetCurrentContext() != w || v.HWWindow == nil || v.HWWindow == w {
		t.Fatal("reset did not return to an independent frontend context")
	}

	v.BeginCore()
	var readFBO, texture, pbo uint32
	gl.GenFramebuffers(1, &readFBO)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, readFBO)
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.GenBuffers(1, &pbo)
	gl.BindBuffer(gl.PIXEL_PACK_BUFFER, pbo)
	gl.PixelStorei(gl.PACK_ROW_LENGTH, 256)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.BlendEquation(gl.FUNC_REVERSE_SUBTRACT)
	v.EndCore()
	v.BeginFrontendFrame()
	checkInt := func(name uint32, want int32) {
		t.Helper()
		var got int32
		gl.GetIntegerv(name, &got)
		if got != want {
			t.Errorf("GL state %#x: got %d, want %d", name, got, want)
		}
	}
	checkInt(gl.PIXEL_PACK_BUFFER_BINDING, 0)
	checkInt(gl.PACK_ROW_LENGTH, 0)
	checkInt(gl.READ_FRAMEBUFFER_BINDING, 0)
	checkInt(gl.BLEND_EQUATION_RGB, gl.FUNC_ADD)
	v.BeginCore()
	checkInt(gl.PIXEL_PACK_BUFFER_BINDING, int32(pbo))
	checkInt(gl.PACK_ROW_LENGTH, 256)
	checkInt(gl.READ_FRAMEBUFFER_BINDING, int32(readFBO))
	checkInt(gl.DRAW_FRAMEBUFFER_BINDING, int32(v.fboID))
	checkInt(gl.BLEND_EQUATION_RGB, gl.FUNC_REVERSE_SUBTRACT)
	checkInt(gl.ACTIVE_TEXTURE, gl.TEXTURE1)
	checkInt(gl.TEXTURE_BINDING_2D, int32(texture))
	if !gl.IsEnabled(gl.BLEND) || !gl.IsEnabled(gl.SCISSOR_TEST) {
		t.Error("frontend rendering changed core enables")
	}
	gl.BindBuffer(gl.PIXEL_PACK_BUFFER, 0)
	gl.PixelStorei(gl.PACK_ROW_LENGTH, 0)
	gl.Disable(gl.SCISSOR_TEST)
	gl.Disable(gl.BLEND)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, v.fboID)

	// Verify the produced pixels are visible through the shared texture.
	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	v.EndCore()
	var frontendFBO uint32
	gl.GenFramebuffers(1, &frontendFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frontendFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, v.texID, 0)
	var pixel [4]byte
	gl.ReadPixels(0, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixel[0]))
	if pixel != [4]byte{255, 0, 0, 255} {
		t.Errorf("shared frame: got %v, want red", pixel)
	}
	v.BeginFrontendFrame()

	// Backbuffer snapshots must retain their orientation for either core origin.
	v.defaultProgram, err = newProgram(vertexShader, defaultFragmentShader)
	if err != nil {
		t.Fatal(err)
	}
	gl.GenVertexArrays(1, &v.vao)
	gl.BindVertexArray(v.vao)
	gl.GenBuffers(1, &v.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, v.vbo)
	vert := uint32(gl.GetAttribLocation(v.defaultProgram, gl.Str("vert\x00")))
	uv := uint32(gl.GetAttribLocation(v.defaultProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(vert)
	gl.VertexAttribPointerWithOffset(vert, 2, gl.FLOAT, false, 16, 0)
	gl.EnableVertexAttribArray(uv)
	gl.VertexAttribPointerWithOffset(uv, 2, gl.FLOAT, false, 16, 8)
	v.width, v.height = 64, 64
	for _, bottomLeft := range []bool{false, true} {
		hw.BottomLeftOrigin = bottomLeft
		v.BeginFrontendFrame()
		fbw, fbh := w.GetFramebufferSize()
		gl.ClearColor(1, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Enable(gl.SCISSOR_TEST)
		gl.Scissor(0, int32(fbh/2), int32(fbw), int32(fbh/2))
		gl.ClearColor(0, 0, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Disable(gl.SCISSOR_TEST)
		v.CachePresentedFrame()
		if !v.RenderCachedFrame() {
			t.Fatal("cached frame unavailable")
		}
		gl.ReadPixels(1, 1, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixel[0]))
		if pixel != [4]byte{255, 0, 0, 255} {
			t.Errorf("cached frame, bottom-left=%v: got %v, want red", bottomLeft, pixel)
		}
	}

	hw.ContextDestroy = func() {
		if glfw.GetCurrentContext() != v.HWWindow {
			t.Error("context_destroy ran outside the HW context")
		}
	}
	v.RunContextDestroy()
	v.ResetCoreState()
	if !gl.IsFramebuffer(frontendFBO) {
		t.Error("HW teardown deleted a frontend framebuffer with the same name")
	}
	if v.fboID != 0 || v.HWWindow != nil || v.cachedFrameValid || v.width != 0 {
		t.Error("HW teardown left stale core resources")
	}

	// A subsequent core must get a fresh context, including on the GL 2.x sync fallback.
	if err := v.ResetHWContext(); err != nil {
		t.Fatal(err)
	}
	v.hwSync = false
	v.BeginCore()
	checkInt(gl.PACK_ROW_LENGTH, 0)
	checkInt(gl.PIXEL_PACK_BUFFER_BINDING, 0)
	checkInt(gl.BLEND_EQUATION_RGB, gl.FUNC_ADD)
	if gl.IsEnabled(gl.BLEND) || gl.IsEnabled(gl.SCISSOR_TEST) {
		t.Error("new core inherited the previous core's enables")
	}
	v.EndCore()
	if e := gl.GetError(); e != gl.NO_ERROR {
		t.Errorf("GL error: %#x", e)
	}
}
