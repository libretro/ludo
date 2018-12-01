// Package video takes care on the game display. It also creates the window
// using GLFW. It exports the Refresh callback used by the libretro
// implementation.
package video

import (
	"log"

	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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
	GLVersion uint
	Window    WindowInterface
	Geom      libretro.GameGeometry
	Font      *glfont.Font

	program        uint32 // default program used for the game window
	roundedProgram uint32 // program to draw rectangles with rounded corners
	borderProgram  uint32 // program to draw rectangles borders
	circleProgram  uint32 // program to draw textured circles
	demulProgram   uint32 // program to draw premultiplied alpha images
	vao            uint32
	vbo            uint32
	texID          uint32
	white          uint32
	pitch          int32
	pixFmt         uint32
	pixType        uint32
	bpp            int32
}

// Init instanciates the video package
func Init(fullscreen bool, GLVersion uint) *Video {
	vid := &Video{}
	vid.GLVersion = GLVersion
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

func (video *Video) configureContext() uint {
	var GLSLVersion uint
	switch video.GLVersion {
	case 20:
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
		GLSLVersion = 110
	case 21:
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
		GLSLVersion = 120
	case 30:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
		GLSLVersion = 130
	case 31:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
		GLSLVersion = 140
	case 32:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		GLSLVersion = 150
	case 41:
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		GLSLVersion = 410
	case 42:
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		GLSLVersion = 420
	default: // 32
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		GLSLVersion = 150
	}
	return GLSLVersion
}

// Configure instanciates the video package
func (video *Video) Configure(fullscreen bool) {
	var width, height int
	var m *glfw.Monitor
	GLSLVersion := video.configureContext()

	if fullscreen {
		m = glfw.GetMonitors()[settings.Settings.VideoMonitorIndex]
		vms := m.GetVideoModes()
		vm := vms[len(vms)-1]
		width = vm.Width
		height = vm.Height
	} else {
		width = 320 * 3
		height = 180 * 3
	}

	var err error
	video.Window, err = glfw.CreateWindow(width, height, "Ludo", m, nil)
	if err != nil {
		panic(err)
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
	video.CoreRatioViewport(fbw, fbh)

	// LoadFont (fontfile, font scale, window width, window height)
	video.Font, err = glfont.LoadFont("assets/font.ttf", int32(64), fbw, fbh, GLSLVersion)
	if err != nil {
		panic(err)
	}

	if state.Global.Verbose {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("[Video]: OpenGL version:", version)
	}

	// Configure the vertex and fragment shaders
	video.program, err = newProgram(GLSLVersion, vertexShader, darkenFragmentShader)
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

	gl.UseProgram(video.program)

	textureUniform := gl.GetUniformLocation(video.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(video.program, 0, gl.Str("outputColor\x00"))

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

	// Sets a default pixel format
	if video.pixFmt == 0 {
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
	}

	gl.GenTextures(1, &video.texID)

	gl.ActiveTexture(gl.TEXTURE0)
	if video.texID == 0 && state.Global.Verbose {
		log.Println("[Video]: Failed to create the vid texture")
	}

	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	video.white = newWhite()
}

// SetPixelFormat is a callback passed to the libretro implementation.
// It allows the core or the game to tell us which pixel format should be used for the display.
func (video *Video) SetPixelFormat(format uint32) bool {
	if state.Global.Verbose {
		log.Printf("[Video]: Set Pixel Format: %v\n", format)
	}

	switch format {
	case libretro.PixelFormat0RGB1555:
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
	case libretro.PixelFormatXRGB8888:
		video.pixFmt = gl.UNSIGNED_INT_8_8_8_8_REV
		video.pixType = gl.BGRA
		video.bpp = 4
	case libretro.PixelFormatRGB565:
		video.pixFmt = gl.UNSIGNED_SHORT_5_6_5
		video.pixType = gl.RGB
		video.bpp = 2
	default:
		log.Fatalf("Unknown pixel type %v", format)
	}

	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)

	return true
}

func (video *Video) updateMaskUniform() {
	maskUniform := gl.GetUniformLocation(video.program, gl.Str("mask\x00"))
	if state.Global.MenuActive {
		gl.Uniform1f(maskUniform, 1.0)
	} else {
		gl.Uniform1f(maskUniform, 0.0)
	}
}

// CoreRatioViewport configures the vertex array to display the game at the center of the window
// while preserving the original ascpect ratio of the game or core
func (video *Video) CoreRatioViewport(fbWidth int, fbHeight int) {
	// Scale the content to fit in the viewport.
	fbw := float32(fbWidth)
	fbh := float32(fbHeight)
	w := fbw
	h := fbw / float32(video.Geom.AspectRatio)
	if h > fbh {
		w = fbh * float32(video.Geom.AspectRatio)
		h = fbh
	}

	// Place the content in the middle of the window.
	x := (fbw - w) / 2
	y := (fbh - h) / 2

	x1, y1, x2, y2, x3, y3, x4, y4 := XYWHTo4points(x, y, w, h, fbh)

	va := []float32{
		//  X, Y, U, V
		x1/fbw*2 - 1, y1/fbh*2 - 1, 0, 1, // left-bottom
		x2/fbw*2 - 1, y2/fbh*2 - 1, 0, 0, // left-top
		x3/fbw*2 - 1, y3/fbh*2 - 1, 1, 1, // right-bottom
		x4/fbw*2 - 1, y4/fbh*2 - 1, 1, 0, // right-top
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
}

// FullViewport configures the GL viewport to take all the available space in the window
func (video *Video) FullViewport() {
	w, h := video.Window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))
}

// Render the current frame
func (video *Video) Render() {
	if state.Global.CoreRunning {
		gl.ClearColor(0, 0, 0, 1)
	} else {
		gl.ClearColor(1, 1, 1, 1)
	}
	gl.Clear(gl.COLOR_BUFFER_BIT)

	fbw, fbh := video.Window.GetFramebufferSize()
	video.CoreRatioViewport(fbw, fbh)

	gl.UseProgram(video.program)
	video.updateMaskUniform()
	gl.Uniform4f(gl.GetUniformLocation(video.program, gl.Str("texColor\x00")), 1, 1, 1, 1)

	gl.BindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// Refresh the texture framebuffer
func (video *Video) Refresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, video.pixType, video.pixFmt, nil)

	video.pitch = pitch
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)

	if data != nil {
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, width, height, video.pixType, video.pixFmt, data)
	}
}

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
