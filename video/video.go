// Package video takes care on the game display. It also creates the window
// using GLFW. It exports the Refresh callback used by the libretro
// implementation.
package video

import (
	"log"
	"path/filepath"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

// Video holds the state of the video package
type Video struct {
	Window   *glfw.Window
	HWWindow *glfw.Window
	Geom     libretro.GameGeometry
	Font     *Font
	FontSm   *Font
	FontLg   *Font
	title    string

	program              uint32 // current program used for the game quad
	defaultProgram       uint32 // default program used for the game quad
	sharpBilinearProgram uint32 // sharp bilinear program used for the game quad
	zfastCRTProgram      uint32 // fast CRT program used for the game quad
	zfastLCDProgram      uint32 // fast LCD program used for the game quad
	roundedProgram       uint32 // program to draw rectangles with rounded corners
	borderProgram        uint32 // program to draw rectangles borders
	circleProgram        uint32 // program to draw textured circles
	vao                  uint32 // vertex array object
	vbo                  uint32 // vertex buffer object
	texID                uint32 // texture id
	fboID                uint32 // framebuffer id for hw-render cores
	rboID                uint32 // depth/stencil renderbuffer for hw-render cores
	identityMat          mgl32.Mat4
	orthoMat             mgl32.Mat4
	texWidth             int32
	texHeight            int32

	pitch         int32  // pitch set by the refresh callback
	pixFmt        uint32 // format set by the environment callback
	pixType       uint32 // pixel type for the core framebuffer
	bpp           int32  // bit per pixel for the core framebuffer
	width, height int32  // dimensions set by the refresh callback
	rot           uint   // rotation index

	needUpload bool // true when the texture needs to be uploaded to the GPU
	data       unsafe.Pointer

	hwContextType       uint32
	hwContextMajor      int
	hwContextMinor      int
	hwDebugContext      bool
	hwContextConfigured bool
}

// Init instanciates the video package
func Init(fullscreen bool) *Video {
	vid := &Video{title: "Ludo"}
	vid.identityMat = mgl32.Ident4()
	vid.Configure(fullscreen)
	return vid
}

// Reconfigure destroys and recreates the window with new attributes
func (video *Video) Reconfigure(fullscreen bool) {
	if video.Window != nil {
		if video.HWWindow != nil {
			video.HWWindow.Destroy()
			video.HWWindow = nil
		}
		// This is the expected frontend behavior and Flycast requires this
		// for fullscreen toggling to work, but ppsspp breaks. OTOH, ppsspp
		// breaks in those situations even if we don't call context_destroy
		// so ignore it.
		hw := state.Core.HWRenderCallback
		if state.CoreRunning && hw != nil && hw.ContextDestroy != nil {
			state.Core.HWRenderCallback.ContextDestroy()
		}
		video.Window.Destroy()
	}
	video.Configure(fullscreen)
}

// CreateSharedHWContext creates a hidden context shared with the visible frontend context.
func (video *Video) CreateSharedHWContext() error {
	if video.HWWindow != nil || video.Window == nil || !state.CoreSetSharedContext {
		return nil
	}

	video.configureWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.False)

	width, height := video.Window.GetSize()
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	hwWindow, err := glfw.CreateWindow(width, height, "", nil, video.Window)
	if err != nil {
		return err
	}

	video.HWWindow = hwWindow
	video.MakeFrontendContextCurrent()
	return nil
}

// DestroySharedHWContext destroys the hidden shared context if it exists.
func (video *Video) DestroySharedHWContext() {
	if video.HWWindow != nil {
		video.HWWindow.Destroy()
		video.HWWindow = nil
	}
}

// MakeFrontendContextCurrent switches back to the visible frontend context.
func (video *Video) MakeFrontendContextCurrent() {
	if video.Window != nil {
		video.Window.MakeContextCurrent()
	}
}

// MakeHardwareContextCurrent switches to the shared HW context when available.
func (video *Video) MakeHardwareContextCurrent() {
	if video.HWWindow != nil {
		video.HWWindow.MakeContextCurrent()
	}
}

// BeginCoreFrame makes the hardware context current before the core issues GL calls.
func (video *Video) BeginCoreFrame() {
	video.MakeHardwareContextCurrent()
}

// EndCoreFrame flushes HW GL work and restores the visible frontend context.
func (video *Video) EndCoreFrame() {
	if video.HWWindow != nil {
		gl.Flush()
		video.MakeFrontendContextCurrent()
	}
}

// GetFramebufferSize retrieves the size, in pixels, of the framebuffer of the specified window.
func (video *Video) GetFramebufferSize() (int, int) {
	if video.Window == nil {
		return 0, 0
	}
	return video.Window.GetFramebufferSize()
}

// SetTitle sets the window title, encoded as UTF-8, of the window.
func (video *Video) SetTitle(title string) {
	video.title = title
	if video.Window == nil {
		return
	}
	video.Window.SetTitle(title)
}

// SetHWRenderContext stores the context requested by the current libretro core.
func (video *Video) SetHWRenderContext(hw *libretro.HWRenderCallback) {
	if hw == nil {
		video.hwContextType = 0
		video.hwContextMajor = 0
		video.hwContextMinor = 0
		video.hwDebugContext = false
		video.hwContextConfigured = false
		return
	}

	video.hwContextType = hw.HWContextType
	video.hwContextMajor = int(hw.VersionMajor)
	video.hwContextMinor = int(hw.VersionMinor)
	video.hwDebugContext = hw.DebugContext
	video.hwContextConfigured = true
}

func (video *Video) configureWindowHints() {
	glfw.DefaultWindowHints()

	if !video.hwContextConfigured {
		return
	}

	switch video.hwContextType {
	case libretro.HWContextOpenGLCore:
		major := video.hwContextMajor
		minor := video.hwContextMinor
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		if major > 0 && !(major == 3 && minor == 1) {
			glfw.WindowHint(glfw.ContextVersionMajor, major)
			glfw.WindowHint(glfw.ContextVersionMinor, minor)
			if major > 3 || (major == 3 && minor >= 2) {
				glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			}
		}
	case libretro.HWContextOpenGL:
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
		if video.hwContextMajor > 0 {
			glfw.WindowHint(glfw.ContextVersionMajor, video.hwContextMajor)
			glfw.WindowHint(glfw.ContextVersionMinor, video.hwContextMinor)
			if video.hwContextMajor > 3 || (video.hwContextMajor == 3 && video.hwContextMinor >= 2) {
				glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)
			}
		}
	case libretro.HWContextOpenGLES2:
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
	case libretro.HWContextOpenGLES3:
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
	case libretro.HWContextOpenGLESVersion:
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
		if video.hwContextMajor > 0 {
			glfw.WindowHint(glfw.ContextVersionMajor, video.hwContextMajor)
			glfw.WindowHint(glfw.ContextVersionMinor, video.hwContextMinor)
		}
	}

	if video.hwDebugContext {
		glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	}
}

// SetShouldClose sets the value of the close flag of the window.
func (video *Video) SetShouldClose(b bool) {
	if video.Window == nil {
		return
	}
	video.Window.SetShouldClose(b)
}

func panicOnErr(v uint32, err error) uint32 {
	if err != nil {
		panic(err)
	}
	return v
}

// Configure instanciates the video package
func (video *Video) Configure(fullscreen bool) {
	var width, height int
	var m *glfw.Monitor

	if fullscreen {
		m = glfw.GetMonitors()[settings.Current.VideoMonitorIndex]
		vms := m.GetVideoModes()
		vm := vms[len(vms)-1]
		width = vm.Width
		height = vm.Height
	} else {
		width = 384 * 2
		height = 240 * 2
	}

	video.configureWindowHints()

	var err error
	video.Window, err = glfw.CreateWindow(width, height, video.title, m, nil)
	if err != nil {
		panic("Window creation failed:" + err.Error())
	}

	video.Window.MakeContextCurrent()

	// Force a minimum size for the window.
	video.Window.SetSizeLimits(160, 120, glfw.DontCare, glfw.DontCare)

	video.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	fbw, fbh := video.Window.GetFramebufferSize()

	// LoadFont (fontfile, font scale, window width, window height)
	fontPath := filepath.Join(settings.Current.AssetsDirectory, "font.ttf")
	video.Font, err = LoadFont(fontPath, int32(36), fbw, fbh)
	if err != nil {
		panic(err)
	}
	video.FontSm, err = LoadFont(fontPath, int32(29), fbw, fbh)
	if err != nil {
		panic(err)
	}
	video.FontLg, err = LoadFont(fontPath, int32(50), fbw, fbh)
	if err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shader
	video.defaultProgram = panicOnErr(newProgram(vertexShader, defaultFragmentShader))
	video.sharpBilinearProgram = panicOnErr(newProgram(vertexShader, sharpBilinearFragmentShader))
	video.zfastCRTProgram = panicOnErr(newProgram(vertexShader, zfastCRTFragmentShader))
	video.zfastLCDProgram = panicOnErr(newProgram(vertexShader, zfastLCDFragmentShader))
	video.roundedProgram = panicOnErr(newProgram(vertexShader, roundedFragmentShader))
	video.borderProgram = panicOnErr(newProgram(vertexShader, borderFragmentShader))
	video.circleProgram = panicOnErr(newProgram(vertexShader, circleFragmentShader))

	video.UpdateFilter(settings.Current.VideoFilter)

	gl.UseProgram(video.program)
	textureUniform := gl.GetUniformLocation(video.program, gl.Str("Texture\x00"))
	gl.Uniform1i(textureUniform, 0)

	// Configure the vertex data
	genVertexArrays(1, &video.vao)
	bindVertexArray(video.vao)

	gl.GenBuffers(1, &video.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(video.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(video.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)

	// Some cores won't call SetPixelFormat, provide default values
	if video.pixFmt == 0 {
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
	}

	if video.Geom.MaxWidth == 0 || video.Geom.MaxHeight == 0 {
		video.Geom.MaxWidth = video.Geom.BaseWidth
		video.Geom.MaxHeight = video.Geom.BaseHeight
	}

	gl.GenTextures(1, &video.texID)
	if video.texID == 0 && state.Verbose {
		log.Fatalln("[Video]: Failed to create the vid texture")
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	video.texWidth = int32(video.Geom.MaxWidth)
	video.texHeight = int32(video.Geom.MaxHeight)
	if video.texWidth < 1 {
		video.texWidth = 1
	}
	if video.texHeight < 1 {
		video.texHeight = 1
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, video.texWidth, video.texHeight, 0, video.pixType, video.pixFmt, nil)

	video.UpdateFilter(settings.Current.VideoFilter)

	video.coreRatioViewport(fbw, fbh, video.Geom.BaseWidth, video.Geom.BaseHeight)

	bindVertexArray(0)

	for e := gl.GetError(); e != gl.NO_ERROR; e = gl.NO_ERROR {
		log.Printf("[Video] OpenGL error: %d\n", e)
	}
}

// UpdateFilter configures the game texture filter and shader. We currently
// support 4 modes:
// Raw: nearest
// Smooth: linear
// Pixel Perfect: sharp-bilinear
// CRT: zfast-crt
// LCD: zfast-lcd
func (video *Video) UpdateFilter(filter string) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	switch filter {
	case "Smooth":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.defaultProgram
	case "Pixel Perfect":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.sharpBilinearProgram
	case "CRT":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.zfastCRTProgram
	case "LCD":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.zfastLCDProgram
	case "Raw":
		fallthrough
	default:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		video.program = video.defaultProgram
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(video.texWidth), float32(video.texHeight))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(video.width), float32(video.height))
	gl.UseProgram(0)
}

// SetPixelFormat is a callback passed to the libretro implementation.
// It allows the core or the game to tell us which pixel format should be used for the display.
func (video *Video) SetPixelFormat(format uint32) bool {
	if state.Verbose {
		log.Printf("[Video]: Set Pixel Format: %v\n", format)
	}

	// PixelStorei also needs to be updated whenever bpp changes
	defer func() { video.needUpload = true }()

	switch format {
	case libretro.PixelFormat0RGB1555:
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
		return true
	case libretro.PixelFormatXRGB8888:
		video.pixFmt = gl.UNSIGNED_INT_8_8_8_8_REV
		video.pixType = gl.BGRA
		video.bpp = 4
		return true
	case libretro.PixelFormatRGB565:
		video.pixFmt = gl.UNSIGNED_SHORT_5_6_5
		video.pixType = gl.RGB
		video.bpp = 2
		return true
	default:
		log.Printf("Unknown pixel type %v", format)
	}

	return false
}

// ResetPitch should be called when unloading a game so that the next game won't
// be rendered with the wrong pitch
func (video *Video) ResetPitch() {
	video.pitch = 0
}

// ResetRot should be called when unloading a game so that the next game won't
// be rendered with the wrong rotation
func (video *Video) ResetRot() {
	video.rot = 0
}

// coreRatioViewport configures the vertex array to display the game at the center of the window
// while preserving the original ascpect ratio of the game or core
func (video *Video) coreRatioViewport(fbWidth, fbHeight, clipWidth, clipHeight int) (x, y, w, h float32) {
	// Scale the content to fit in the viewport.
	fbw := float32(fbWidth)
	fbh := float32(fbHeight)

	// NXEngine workaround
	aspectRatio := float32(video.Geom.AspectRatio)
	if aspectRatio == 0 {
		aspectRatio = float32(video.Geom.BaseWidth) / float32(video.Geom.BaseHeight)
	}

	h = fbh
	w = fbh * aspectRatio
	if w > fbw {
		h = fbw / aspectRatio
		w = fbw
	}

	// Place the content in the middle of the window.
	x = (fbw - w) / 2
	y = (fbh - h) / 2

	va := video.vertexArray(x, y, w, h, 1.0)

	if clipWidth == 0 {
		clipWidth = int(video.texWidth)
	}
	if clipHeight == 0 {
		clipHeight = int(video.texHeight)
	}

	va[3] = float32(clipHeight) / float32(video.texHeight)
	va[10] = float32(clipWidth) / float32(video.texWidth)
	va[11] = va[3]
	va[14] = va[10]

	va = rotateUV(va, video.rot)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	return
}

// PrepareCoreContext resets frontend GL state so HW cores see a minimal context
func (video *Video) PrepareCoreContext() {
	gl.UseProgram(0)
	bindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	bindBackbuffer()
	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.DITHER)
	gl.Disable(gl.SCISSOR_TEST)
	gl.Disable(gl.STENCIL_TEST)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.PixelStorei(gl.PACK_ROW_LENGTH, 0)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 4)
	gl.PixelStorei(gl.PACK_ALIGNMENT, 4)
}

// ResizeViewport resizes the GL viewport to the framebuffer size
func (video *Video) ResizeViewport() {
	fbw, fbh := video.Window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbw), int32(fbh))
	if video.Font != nil {
		video.Font.UpdateResolution(fbw, fbh)
	}
	if video.FontSm != nil {
		video.FontSm.UpdateResolution(fbw, fbh)
	}
	if video.FontLg != nil {
		video.FontLg.UpdateResolution(fbw, fbh)
	}
}

// Render the current frame
func (video *Video) Render() {
	// Render directly to the screen
	bindBackbuffer()

	// We can't trust the core to leave the OpenGL in the same state as
	// before retro_run() was called so we restore some state manually.
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DITHER)
	gl.Disable(gl.STENCIL_TEST)
	gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.Enable(gl.TEXTURE_2D)

	video.ResizeViewport()

	if !state.CoreRunning {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		return
	}

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Early return to not render the first frame of a newly loaded game with the
	// previous game pitch. A sane pitch must be set by video.Refresh first.
	if state.Core.HWRenderCallback == nil && video.pitch == 0 {
		return
	}

	video.uploadTexture()

	fbw, fbh := video.Window.GetFramebufferSize()
	_, _, w, h := video.coreRatioViewport(fbw, fbh, int(video.width), int(video.height))

	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("OutputSize\x00")), w, h)

	if state.Core.HWRenderCallback != nil {
		gl.UniformMatrix4fv(gl.GetUniformLocation(video.program, gl.Str("MVP\x00")), 1, false, &video.orthoMat[0])
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	bindVertexArray(video.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	bindVertexArray(0)

	// Reset MVP to identity to avoid menu issues
	gl.UniformMatrix4fv(gl.GetUniformLocation(video.program, gl.Str("MVP\x00")), 1, false, &video.identityMat[0])
	gl.UseProgram(0)
}

// Refresh the texture framebuffer
func (video *Video) Refresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	video.needUpload = true
	video.width = width
	video.height = height
	video.pitch = pitch
	video.data = data
}

func (video *Video) uploadTexture() {
	if !video.needUpload || video.data == nil {
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	if video.data != libretro.HWFrameBufferValid {
		video.texWidth = int32(video.Geom.MaxWidth)
		video.texHeight = int32(video.Geom.MaxHeight)
		if video.texWidth == 0 {
			video.texWidth = video.width
		}
		if video.texHeight == 0 {
			video.texHeight = video.height
		}
		if video.texWidth < 1 {
			video.texWidth = 1
		}
		if video.texHeight < 1 {
			video.texHeight = 1
		}
		gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, video.texWidth, video.texHeight, 0, video.pixType, video.pixFmt, video.data)
	}

	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(video.texWidth), float32(video.texHeight))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(video.width), float32(video.height))
	gl.UseProgram(0)
	video.needUpload = false
}

// CurrentFramebuffer returns the current FBO ID
func (video *Video) CurrentFramebuffer() uintptr {
	return uintptr(video.fboID)
}

// ProcAddress returns the address of the proc from GLFW
func (video *Video) ProcAddress(procName string) uintptr {
	return uintptr(glfw.GetProcAddress(procName))
}

// SetRotation rotates the game image as requested by the core
func (video *Video) SetRotation(rot uint) bool {
	// limit to valid values (0, 1, 2, 3, which rotates screen by 0, 90, 180 270 degrees counter-clockwise)
	video.rot = rot % 4

	if state.Verbose {
		log.Printf("[Video]: Set Rotation: %v", video.rot)
	}

	return true
}

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
