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
	fboID                uint32
	rboID                uint32
	identityMat          [16]float32
	orthoMat             [16]float32
	texWidth             int32
	texHeight            int32

	pitch         int32  // pitch set by the refresh callback
	pixFmt        uint32 // format set by the environment callback
	pixType       uint32 // pixel type for the core framebuffer
	bpp           int32  // bit per pixel for the core framebuffer
	width, height int32  // dimensions set by the refresh callback
	rot           uint   // rotation index

	hwContextType       uint32
	hwContextMajor      int
	hwContextMinor      int
	hwDebugContext      bool
	hwContextConfigured bool
	presentDirty        bool
	presentFBWidth      int
	presentFBHeight     int
	presentClipWidth    int32
	presentClipHeight   int32
	presentRot          uint
}

// Init instanciates the video package
func Init(fullscreen bool) *Video {
	vid := &Video{
		identityMat: identityMatrix(),
		orthoMat:    identityMatrix(),
		title:       "Ludo",
	}
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
		video.Window.Destroy()
	}
	video.resetFrontendState()
	video.Configure(fullscreen)
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

// SetHWRenderContext stores the context requested by the current libretro core
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

func (video *Video) textureSize() (int32, int32) {
	if video.texWidth > 0 && video.texHeight > 0 {
		return video.texWidth, video.texHeight
	}

	w := int32(video.Geom.MaxWidth)
	h := int32(video.Geom.MaxHeight)
	if w == 0 {
		if video.width > 0 {
			w = video.width
		} else {
			w = int32(video.Geom.BaseWidth)
		}
	}
	if h == 0 {
		if video.height > 0 {
			h = video.height
		} else {
			h = int32(video.Geom.BaseHeight)
		}
	}
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return w, h
}

func (video *Video) resetFrontendState() {
	video.Window = nil
	video.HWWindow = nil
	video.Font = nil
	video.FontSm = nil
	video.FontLg = nil
	video.program = 0
	video.defaultProgram = 0
	video.sharpBilinearProgram = 0
	video.zfastCRTProgram = 0
	video.zfastLCDProgram = 0
	video.roundedProgram = 0
	video.borderProgram = 0
	video.circleProgram = 0
	video.vao = 0
	video.vbo = 0
	video.texID = 0
	video.fboID = 0
	video.rboID = 0
	video.presentDirty = true
	video.presentFBWidth = 0
	video.presentFBHeight = 0
	video.presentClipWidth = 0
	video.presentClipHeight = 0
	video.presentRot = 0
	libretro.SetCurrentFramebufferValue(0)
}

func (video *Video) createSharedHWContext() error {
	if video.Window == nil || state.Core == nil || state.Core.HWRenderCallback == nil || state.CoreRunning || !state.CoreSetSharedContext {
		return nil
	}

	video.configureWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.False)

	hwWindow, err := glfw.CreateWindow(1, 1, "", nil, video.Window)
	if err != nil {
		return err
	}

	video.HWWindow = hwWindow
	video.MakeFrontendContextCurrent()
	return nil
}

// MakeFrontendContextCurrent switches back to the visible window context
func (video *Video) MakeFrontendContextCurrent() {
	if video.HWWindow != nil && video.Window != nil {
		video.Window.MakeContextCurrent()
	}
}

// MakeHardwareContextCurrent switches to the shared HW context when available
func (video *Video) MakeHardwareContextCurrent() {
	if video.HWWindow != nil {
		video.HWWindow.MakeContextCurrent()
	}
}

// BeginCoreFrame makes the hardware context current before the core issues GL calls
func (video *Video) BeginCoreFrame() {
	video.MakeHardwareContextCurrent()
}

// EndCoreFrame flushes HW GL work and restores the visible frontend context
func (video *Video) EndCoreFrame() {
	if video.HWWindow != nil {
		gl.Flush()
	}
	video.MakeFrontendContextCurrent()
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

func (video *Video) invalidatePresentGeometry() {
	video.presentDirty = true
}

func (video *Video) syncPresentGeometry(clipWidth, clipHeight int32) (float32, float32) {
	fbw, fbh := video.Window.GetFramebufferSize()
	if video.presentDirty ||
		video.presentFBWidth != fbw ||
		video.presentFBHeight != fbh ||
		video.presentClipWidth != clipWidth ||
		video.presentClipHeight != clipHeight ||
		video.presentRot != video.rot {
		_, _, w, h := video.coreRatioViewport(fbw, fbh, int(clipWidth), int(clipHeight))
		video.presentFBWidth = fbw
		video.presentFBHeight = fbh
		video.presentClipWidth = clipWidth
		video.presentClipHeight = clipHeight
		video.presentRot = video.rot
		video.presentDirty = false
		return w, h
	}

	aspectRatio := float32(video.Geom.AspectRatio)
	if aspectRatio == 0 && clipWidth > 0 && clipHeight > 0 {
		aspectRatio = float32(clipWidth) / float32(clipHeight)
	}
	if aspectRatio == 0 {
		aspectRatio = 1
	}
	w := float32(fbh) * aspectRatio
	h := float32(fbh)
	if w > float32(fbw) {
		h = float32(fbw) / aspectRatio
		w = float32(fbw)
	}
	return w, h
}

func (video *Video) presentFrame() {
	if video.Window == nil || video.vao == 0 || video.program == 0 {
		return
	}

	video.MakeFrontendContextCurrent()
	bindBackbuffer()

	fbw, fbh := video.Window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbw), int32(fbh))
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DITHER)
	gl.Disable(gl.STENCIL_TEST)
	gl.Disable(gl.BLEND)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	w, h := video.syncPresentGeometry(video.width, video.height)

	gl.UseProgram(video.program)
	video.setGameProgramColor(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("OutputSize\x00")), w, h)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(video.texWidth), float32(video.texHeight))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(video.width), float32(video.height))

	mvp := video.identityMat
	if state.Core.HWRenderCallback != nil {
		mvp = video.orthoMat
	}
	gl.UniformMatrix4fv(gl.GetUniformLocation(video.program, gl.Str("MVP\x00")), 1, false, &mvp[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	bindVertexArray(video.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	bindVertexArray(0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(video.program, gl.Str("MVP\x00")), 1, false, &video.identityMat[0])
	gl.UseProgram(0)
}

func (video *Video) initSharedTexture() {
	gl.GenTextures(1, &video.texID)
	if video.texID == 0 && state.Verbose {
		log.Println("[Video]: Failed to create the vid texture")
	}
	if video.pixFmt == 0 {
		// Some cores won't call SetPixelFormat, provide default values
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
	}
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (video *Video) InitPresentResources() {
	if video.Window == nil || video.defaultProgram != 0 {
		return
	}

	video.defaultProgram = panicOnErr(newProgram(vertexShader, defaultFragmentShader))
	if video.sharpBilinearProgram == 0 {
		video.sharpBilinearProgram = panicOnErr(newProgram(vertexShader, sharpBilinearFragmentShader))
	}
	if video.zfastCRTProgram == 0 {
		video.zfastCRTProgram = panicOnErr(newProgram(vertexShader, zfastCRTFragmentShader))
	}
	if video.zfastLCDProgram == 0 {
		video.zfastLCDProgram = panicOnErr(newProgram(vertexShader, zfastLCDFragmentShader))
	}
	video.program = video.defaultProgram

	genVertexArrays(1, &video.vao)
	bindVertexArray(video.vao)

	gl.GenBuffers(1, &video.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*4, 2*4)

	fbw, fbh := video.Window.GetFramebufferSize()
	video.coreRatioViewport(fbw, fbh, video.Geom.BaseWidth, video.Geom.BaseHeight)

	bindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.UseProgram(0)

	for e := gl.GetError(); e != gl.NO_ERROR; e = gl.GetError() {
		log.Printf("[Video] OpenGL error: %d\n", e)
	}

	video.UpdateFilter(settings.Current.VideoFilter)
}

func (video *Video) InitUIResources() {
	if video.Window == nil {
		return
	}

	video.InitPresentResources()

	var err error
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

	video.roundedProgram = panicOnErr(newProgram(vertexShader, roundedFragmentShader))
	video.borderProgram = panicOnErr(newProgram(vertexShader, borderFragmentShader))
	video.circleProgram = panicOnErr(newProgram(vertexShader, circleFragmentShader))
}

// Configure instanciates the video package
func (video *Video) Configure(fullscreen bool) {
	var width, height int
	var monitor *glfw.Monitor

	if fullscreen {
		monitor = glfw.GetMonitors()[settings.Current.VideoMonitorIndex]
		vms := monitor.GetVideoModes()
		vm := vms[len(vms)-1]
		width = vm.Width
		height = vm.Height
	} else {
		width = 384 * 2
		height = 240 * 2
	}

	video.configureWindowHints()

	var err error
	video.Window, err = glfw.CreateWindow(width, height, video.title, monitor, nil)
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

	video.initSharedTexture()
	if err := video.createSharedHWContext(); err != nil {
		panic("Shared HW context creation failed:" + err.Error())
	}
	video.InitPresentResources()
	video.InitUIResources()
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
		if video.sharpBilinearProgram != 0 {
			video.program = video.sharpBilinearProgram
		} else {
			video.program = video.defaultProgram
		}
	case "CRT":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		if video.zfastCRTProgram != 0 {
			video.program = video.zfastCRTProgram
		} else {
			video.program = video.defaultProgram
		}
	case "LCD":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		if video.zfastLCDProgram != 0 {
			video.program = video.zfastLCDProgram
		} else {
			video.program = video.defaultProgram
		}
	case "Raw":
		fallthrough
	default:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		video.program = video.defaultProgram
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	textureWidth, textureHeight := video.textureSize()
	gl.UseProgram(video.program)
	gl.Uniform1i(gl.GetUniformLocation(video.program, gl.Str("Texture\x00")), 0)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(textureWidth), float32(textureHeight))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(video.width), float32(video.height))
	gl.UseProgram(0)
	video.invalidatePresentGeometry()
}

func (video *Video) setGameProgramColor(program uint32) {
	loc := gl.GetUniformLocation(program, gl.Str("color\x00"))
	if loc != -1 {
		gl.Uniform4f(loc, 1, 1, 1, 1)
	}
}

// SetPixelFormat is a callback passed to the libretro implementation.
// It allows the core or the game to tell us which pixel format should be used for the display.
func (video *Video) SetPixelFormat(format uint32) bool {
	if state.Verbose {
		log.Printf("[Video]: Set Pixel Format: %v\n", format)
	}

	switch format {
	case libretro.PixelFormat0RGB1555:
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
		video.invalidatePresentGeometry()
		return true
	case libretro.PixelFormatXRGB8888:
		video.pixFmt = gl.UNSIGNED_INT_8_8_8_8_REV
		video.pixType = gl.BGRA
		video.bpp = 4
		video.invalidatePresentGeometry()
		return true
	case libretro.PixelFormatRGB565:
		video.pixFmt = gl.UNSIGNED_SHORT_5_6_5
		video.pixType = gl.RGB
		video.bpp = 2
		video.invalidatePresentGeometry()
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
	video.invalidatePresentGeometry()
}

// coreRatioViewport configures the vertex array to display the game at the center of the window
// while preserving the original ascpect ratio of the game or core
func (video *Video) coreRatioViewport(fbWidth, fbHeight, clipWidth, clipHeight int) (x, y, w, h float32) {
	fbw := float32(fbWidth)
	fbh := float32(fbHeight)

	aspectRatio := float32(video.Geom.AspectRatio)
	if aspectRatio == 0 {
		baseWidth := video.Geom.BaseWidth
		baseHeight := video.Geom.BaseHeight
		if baseWidth == 0 {
			baseWidth = clipWidth
		}
		if baseHeight == 0 {
			baseHeight = clipHeight
		}
		if baseWidth > 0 && baseHeight > 0 {
			aspectRatio = float32(baseWidth) / float32(baseHeight)
		} else {
			aspectRatio = 1
		}
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

	textureWidth, textureHeight := video.textureSize()
	if clipWidth == 0 {
		clipWidth = int(textureWidth)
	}
	if clipHeight == 0 {
		clipHeight = int(textureHeight)
	}

	va := video.vertexArray(x, y, w, h, 1.0)
	va[3] = float32(clipHeight) / float32(textureHeight)
	va[10] = float32(clipWidth) / float32(textureWidth)
	va[11] = va[3]
	va[14] = va[10]

	va = rotateUV(va, video.rot)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	return
}

// ResizeViewport resizes the GL viewport to the framebuffer size
func (video *Video) ResizeViewport() {
	video.MakeFrontendContextCurrent()
	fbw, fbh := video.Window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbw), int32(fbh))
	if video.presentFBWidth != fbw || video.presentFBHeight != fbh {
		video.invalidatePresentGeometry()
	}
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
	video.MakeFrontendContextCurrent()
	bindBackbuffer()

	// We can't trust the core to leave OpenGL state untouched after retro_run().
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DITHER)
	gl.Disable(gl.STENCIL_TEST)
	gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

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

	if state.Core.HWRenderCallback != nil {
		video.presentFrame()
		return
	}

	fbw, fbh := video.Window.GetFramebufferSize()
	_, _, w, h := video.coreRatioViewport(fbw, fbh, int(video.width), int(video.height))

	gl.UseProgram(video.program)
	video.setGameProgramColor(video.program)
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

	gl.UniformMatrix4fv(gl.GetUniformLocation(video.program, gl.Str("MVP\x00")), 1, false, &video.identityMat[0])
	gl.UseProgram(0)
}

// Refresh the texture framebuffer
func (video *Video) Refresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	bindBackbuffer()

	if video.width != width || video.height != height {
		video.invalidatePresentGeometry()
	}
	video.width = width
	video.height = height
	video.pitch = pitch

	textureWidth, textureHeight := video.textureSize()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	if video.texWidth != textureWidth || video.texHeight != textureHeight {
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, textureWidth, textureHeight, 0, video.pixType, video.pixFmt, nil)
		video.texWidth = textureWidth
		video.texHeight = textureHeight
		video.invalidatePresentGeometry()
	}

	if video.bpp > 0 {
		// PixelStorei also needs to be updated whenever bpp changes
		gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)
	}

	if data != nil && data != libretro.HWFrameBufferValid {
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, width, height, video.pixType, video.pixFmt, data)
	}

	if video.bpp > 0 {
		gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	}

	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(textureWidth), float32(textureHeight))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(width), float32(height))
	gl.UseProgram(0)

	if state.Core.HWRenderCallback != nil && video.HWWindow == nil && !state.MenuActive {
		video.presentFrame()
	}
}

// CurrentFramebuffer returns the current FBO ID
func (video *Video) CurrentFramebuffer() uintptr {
	return uintptr(video.fboID)
}

// ProcAddress returns the address of the proc from GLFW
func (video *Video) ProcAddress(procName string) uintptr {
	if glfw.GetCurrentContext() == nil {
		return 0
	}
	addr := glfw.GetProcAddress(procName)
	if addr == nil {
		return 0
	}
	return uintptr(addr)
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
