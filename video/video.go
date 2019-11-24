// Package video takes care on the game display. It also creates the window
// using GLFW. It exports the Refresh callback used by the libretro
// implementation.
package video

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kivutar/glfont"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

// WindowInterface lists all the methods from glfw.Window that we are using.
// It is there only to allow mocking during tests.
type WindowInterface interface {
	GetFramebufferSize() (width, height int)
	Destroy()
	MakeContextCurrent()
	SetSizeLimits(minw, minh, maxw, maxh int)
	SetInputMode(mode glfw.InputMode, value int)
	GetKey(key glfw.Key) glfw.Action
	SetShouldClose(bool)
	ShouldClose() bool
	SetTitle(string)
	SwapBuffers()
}

// Video holds the state of the video package
type Video struct {
	Window WindowInterface
	Geom   libretro.GameGeometry
	Font   *glfont.Font

	program              uint32 // current program used for the game quad
	defaultProgram       uint32 // default program used for the game quad
	sharpBilinearProgram uint32 // sharp bilinear program used for the game quad
	zfastCRTProgram      uint32 // fast CRT program used for the game quad
	roundedProgram       uint32 // program to draw rectangles with rounded corners
	borderProgram        uint32 // program to draw rectangles borders
	circleProgram        uint32 // program to draw textured circles
	demulProgram         uint32 // program to draw premultiplied alpha images
	vao                  uint32
	vbo                  uint32
	texID                uint32
	fboID                uint32
	rboID                uint32

	pitch         int32  // pitch set by the refresh callback
	pixFmt        uint32 // format set by the environment callback
	pixType       uint32
	bpp           int32
	width, height int32 // dimensions set by the refresh callback
}

// Init instanciates the video package
func Init(fullscreen bool) *Video {
	vid := &Video{}
	vid.Configure(fullscreen)
	return vid
}

// Reconfigure destroys and recreates the window with new attributes
func (video *Video) Reconfigure(fullscreen bool) {
	if video.Window != nil {
		video.Window.Destroy()
	}
	video.Configure(fullscreen)
}

func getGLSLVersion() uint {
	GLVersion := gl.GoStr(gl.GetString(gl.VERSION))
	GLSLVersion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	if state.Global.Verbose {
		log.Println("[Video]: OpenGL version:", GLVersion)
		log.Println("[Video]: GLSL version:", GLSLVersion)
	}

	clean := strings.Replace(GLSLVersion[:4], ".", "", -1)
	v, err := strconv.Atoi(clean)
	if err != nil {
		log.Println("[Video]: Couldn't parse GLSL version:", err)
		return 120
	}
	return uint(v)
}

// InitFramebuffer initializes and configures the video frame buffer based on
// informations from the HWRenderCallback of the libretro core.
func (video *Video) InitFramebuffer(width, height int) {
	log.Printf("[Video]: Initializing HW render (%v x %v).\n", width, height)

	gl.GenFramebuffers(1, &video.fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, video.fboID)

	//gl.GenTextures(1, &video.texID)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, int32(width), int32(height))

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, video.texID, 0)

	hw := state.Global.Core.HWRenderCallback

	if hw != nil {
		if hw.Depth && hw.Stencil {
			gl.GenRenderbuffers(1, &video.rboID)
			gl.BindRenderbuffer(gl.RENDERBUFFER, video.rboID)
			gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))

			gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, video.rboID)
		} else if hw.Depth {
			gl.GenRenderbuffers(1, &video.rboID)
			gl.BindRenderbuffer(gl.RENDERBUFFER, video.rboID)
			gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, int32(width), int32(height))

			gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, video.rboID)
		}

		if hw.Depth || hw.Stencil {
			gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
		}
	}

	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalln("[Video] Framebuffer is not complete.")
	}

	gl.ClearColor(0, 0, 0, 1)
	if hw.Depth && hw.Stencil {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	} else if hw.Depth {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	} else {
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Configure instanciates the video package
func (video *Video) Configure(fullscreen bool) {
	var width, height int
	var m *glfw.Monitor

	if fullscreen {
		m = glfw.GetMonitors()[settings.Current.VideoMonitorIndex]
		vm := m.GetVideoMode()
		width = vm.Width
		height = vm.Height
	} else {
		width = 320 * 3
		height = 180 * 3
	}

	// On OSX we have to force a core profile to not end up with 2.1 which cause
	// a font drawing issue
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	} else {
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
	}

	var err error
	video.Window = glfw.CreateWindow(width, height, "Ludo", m, nil)
	if video.Window == nil {
		panic("Window creation failed")
	}

	video.Window.MakeContextCurrent()

	// Force a minimum size for the window.
	video.Window.SetSizeLimits(160, 120, glfw.DontCare, glfw.DontCare)

	video.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	GLSLVersion := getGLSLVersion()

	fbw, fbh := video.Window.GetFramebufferSize()

	// LoadFont (fontfile, font scale, window width, window height)
	assets := settings.Current.AssetsDirectory
	video.Font, err = glfont.LoadFont(assets+"/font.ttf", int32(36*2), fbw, fbh, GLSLVersion)
	if err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shaders
	video.defaultProgram, err = newProgram(GLSLVersion, vertexShader, defaultFragmentShader)
	if err != nil {
		panic(err)
	}

	video.sharpBilinearProgram, err = newProgram(GLSLVersion, vertexShader, sharpBilinearFragmentShader)
	if err != nil {
		panic(err)
	}

	video.zfastCRTProgram, err = newProgram(GLSLVersion, vertexShader, zfastCRTFragmentShader)
	if err != nil {
		panic(err)
	}

	video.roundedProgram, err = newProgram(GLSLVersion, vertexShader, roundedFragmentShader)
	if err != nil {
		panic(err)
	}

	video.borderProgram, err = newProgram(GLSLVersion, vertexShader, borderFragmentShader)
	if err != nil {
		panic(err)
	}

	video.circleProgram, err = newProgram(GLSLVersion, vertexShader, circleFragmentShader)
	if err != nil {
		panic(err)
	}

	video.demulProgram, err = newProgram(GLSLVersion, vertexShader, demulFragmentShader)
	if err != nil {
		panic(err)
	}

	video.UpdateFilter(settings.Current.VideoFilter)

	textureUniform := gl.GetUniformLocation(video.program, gl.Str("Texture\x00"))
	gl.Uniform1i(textureUniform, 0)

	// Configure the vertex data
	gl.GenVertexArrays(1, &video.vao)
	gl.BindVertexArray(video.vao)

	gl.GenBuffers(1, &video.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(video.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(video.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	// Some cores won't call SetPixelFormat, provide default values
	if video.pixFmt == 0 {
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
	}

	gl.GenTextures(1, &video.texID)
	if video.texID == 0 && state.Global.Verbose {
		log.Fatalln("[Video]: Failed to create the vid texture")
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	video.UpdateFilter(settings.Current.VideoFilter)

	video.coreRatioViewport(fbw, fbh)

	if state.Global.CoreRunning && state.Global.Core.HWRenderCallback != nil {
		video.InitFramebuffer(video.Geom.BaseWidth, video.Geom.BaseHeight)
		state.Global.Core.HWRenderCallback.ContextReset()
	}

	if e := gl.GetError(); e != gl.NO_ERROR {
		log.Printf("[Video] OpenGL error: %d\n", e)
	}
}

// UpdateFilter configures the game texture filter and shader. We currently
// support 4 modes: nearest, linear, sharp-bilinear and zfast-crt.
func (video *Video) UpdateFilter(filter string) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	switch filter {
	case "linear":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.defaultProgram
	case "sharp-bilinear":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.sharpBilinearProgram
	case "zfast-crt":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		video.program = video.zfastCRTProgram
	case "nearest":
		fallthrough
	default:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		video.program = video.defaultProgram
	}
	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(video.width), float32(video.height))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(video.width), float32(video.height))
	gl.UseProgram(0)
}

// SetPixelFormat is a callback passed to the libretro implementation.
// It allows the core or the game to tell us which pixel format should be used for the display.
func (video *Video) SetPixelFormat(format uint32) bool {
	if state.Global.Verbose {
		log.Printf("[Video]: Set Pixel Format: %v\n", format)
	}

	// PixelStorei also needs to be updated whenever bpp changes
	defer gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)

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

// coreRatioViewport configures the vertex array to display the game at the center of the window
// while preserving the original ascpect ratio of the game or core
func (video *Video) coreRatioViewport(fbWidth int, fbHeight int) (x, y, w, h float32) {
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
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)

	return
}

// ResizeViewport resizes the GL viewport to the framebuffer size
func (video *Video) ResizeViewport() {
	fbw, fbh := video.Window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbw), int32(fbh))
}

// Render the current frame
func (video *Video) Render() {
	if !state.Global.CoreRunning {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		return
	}
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Early return to not render the first frame of a newly loaded game with the
	// previous game pitch. A sane pitch must be set by video.Refresh first.
	if video.pitch == 0 {
		return
	}

	fbw, fbh := video.Window.GetFramebufferSize()
	_, _, w, h := video.coreRatioViewport(fbw, fbh)

	gl.UseProgram(video.program)
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("OutputSize\x00")), w, h)

	gl.BindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.UseProgram(0)
}

// Refresh the texture framebuffer
func (video *Video) Refresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	video.width = width
	video.height = height
	video.pitch = pitch
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)

	if data != nil && data != libretro.HWFrameBufferValid {
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, width, height, video.pixType, video.pixFmt, data)
	}

	gl.UseProgram(video.program)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, video.pixType, video.pixFmt, nil)

	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("TextureSize\x00")), float32(width), float32(height))
	gl.Uniform2f(gl.GetUniformLocation(video.program, gl.Str("InputSize\x00")), float32(width), float32(height))

	gl.UseProgram(0)
}

// CurrentFramebuffer returns the current FBO ID
func (video *Video) CurrentFramebuffer() uintptr {
	return uintptr(video.fboID)
}

// ProcAddress returns the address of the proc from GLFW
func (video *Video) ProcAddress(procName string) uintptr {
	return uintptr(glfw.GetProcAddress(procName))
}

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
