package main

import (
	"fmt"
	"libretro"
	"log"

	"strings"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/nullboundary/glfont"
)

var window *glfw.Window

var scale = 3.0
var originalAspectRatio float64

var video struct {
	program uint32
	vao     uint32
	vbo     uint32
	texID   uint32
	pitch   int32
	pixFmt  uint32
	pixType uint32
	bpp     int32
	font    *glfont.Font
}

func videoSetPixelFormat(format uint32) bool {
	fmt.Printf("[Video]: Set Pixel Format: %v\n", format)

	if video.texID != 0 {
		log.Fatal("Tried to change pixel format after initialization.")
	}

	switch format {
	case libretro.PixelFormat0RGB1555:
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
		video.pixType = gl.BGRA
		video.bpp = 2
		break
	case libretro.PixelFormatXRGB8888:
		video.pixFmt = gl.UNSIGNED_INT_8_8_8_8_REV
		video.pixType = gl.BGRA
		video.bpp = 4
		break
	case libretro.PixelFormatRGB565:
		video.pixFmt = gl.UNSIGNED_SHORT_5_6_5
		video.pixType = gl.RGB
		video.bpp = 2
		break
	default:
		log.Fatalf("Unknown pixel type %v", format)
	}

	return true
}

func updateMaskUniform() {
	maskUniform := gl.GetUniformLocation(video.program, gl.Str("mask\x00"))
	if menuActive {
		gl.Uniform1f(maskUniform, 1.0)
	} else {
		gl.Uniform1f(maskUniform, 0.0)
	}
}

func createWindow(width int, height int) {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	var err error
	window, err = glfw.CreateWindow(width, height, "playthemall", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Force a minimum size for the window.
	window.SetSizeLimits(160, 120, glfw.DontCare, glfw.DontCare)

	// When resizing the window, also resize the content.
	window.SetFramebufferSizeCallback(resizeFramebuffer)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	//load font (fontfile, font scale, window width, window height
	video.font, err = glfont.LoadFont("font.ttf", int32(64), width, height)
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("[Video]: OpenGL version: ", version)

	// Configure the vertex and fragment shaders
	video.program, err = newProgram(vertexShader, fragmentShader)
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
}

func toggleFullscreen() {
	if window.GetMonitor() == nil {
		m := glfw.GetMonitors()[0]
		vm := m.GetVideoMode()
		newWindow, _ := glfw.CreateWindow(vm.Width, vm.Height, "playthemall", m, window)
		window.Destroy()
		window = newWindow
		window.MakeContextCurrent()
		w, h := window.GetFramebufferSize()
		video.font, _ = glfont.LoadFont("font.ttf", int32(64), w, h)
		gl.Viewport(0, 0, int32(w), int32(h))
	} else {
		avi := core.GetSystemAVInfo()
		geom := avi.Geometry
		nwidth, nheight := resizeToAspect(geom.AspectRatio, float64(geom.BaseWidth), float64(geom.BaseHeight))
		width := int(nwidth * scale)
		height := int(nheight * scale)
		newWindow, _ := glfw.CreateWindow(width, height, "playthemall", nil, window)
		window.Destroy()
		window = newWindow
		window.MakeContextCurrent()
		w, h := window.GetFramebufferSize()
		video.font, _ = glfont.LoadFont("font.ttf", int32(64), w, h)
		gl.Viewport(0, 0, int32(w), int32(h))
		window.SetSizeLimits(160, 120, glfw.DontCare, glfw.DontCare)
		window.SetFramebufferSizeCallback(resizeFramebuffer)
	}
}

func resizeToAspect(ratio float64, sw float64, sh float64) (dw float64, dh float64) {
	if ratio <= 0 {
		ratio = sw / sh
	}

	if sw/sh < 1.0 {
		dw = dh * ratio
		dh = sh
	} else {
		dw = sw
		dh = dw / ratio
	}
	return
}

func resizeFramebuffer(w *glfw.Window, screenWidth int, screenHeight int) {
	// Scale the content to fit in the viewport.
	width := float64(screenWidth)
	height := float64(screenHeight)
	viewWidth := width
	viewHeight := width / originalAspectRatio
	if viewHeight > height {
		viewHeight = height
		viewWidth = height * originalAspectRatio
	}

	// Place the content in the middle of the window.
	vportX := (width - viewWidth) / 2
	vportY := (height - viewHeight) / 2

	gl.Viewport(int32(vportX), int32(vportY), int32(viewWidth), int32(viewHeight))
}

func videoConfigure(geom libretro.GameGeometry) {
	originalAspectRatio = geom.AspectRatio
	nwidth, nheight := resizeToAspect(originalAspectRatio, float64(geom.BaseWidth), float64(geom.BaseHeight))

	width := int(nwidth * scale)
	height := int(nheight * scale)

	if window == nil {
		createWindow(width, height)
	}

	if video.pixFmt == 0 {
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
	}

	gl.GenTextures(1, &video.texID)

	gl.ActiveTexture(gl.TEXTURE0)
	if video.texID == 0 {
		fmt.Println("[Video]: Failed to create the video texture")
	}

	video.pitch = int32(geom.BaseWidth) * video.bpp

	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func renderNotifications() {
	_, height := window.GetSize()
	for i, n := range notifications {
		video.font.SetColor(0.5, 0.5, 0.0, float32(n.frames)/120.0)
		video.font.Printf(30+1, float32(height-30*len(notifications)+30*i+1), 0.25, n.message)
		video.font.SetColor(1.0, 1.0, 0.0, float32(n.frames)/120.0)
		video.font.Printf(30, float32(height-30*len(notifications)+30*i), 0.25, n.message)
	}
}

// Render the current frame
func videoRender() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(video.program)
	updateMaskUniform()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	processNotifications()
	renderNotifications()
}

// Refresh the texture framebuffer
func videoRefresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, video.pixType, video.pixFmt, nil)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	if pitch != video.pitch {
		video.pitch = pitch
		gl.PixelStorei(gl.UNPACK_ROW_LENGTH, video.pitch/video.bpp)
	}

	if data != nil {
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, width, height, video.pixType, video.pixFmt, data)
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 330

in vec2 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(vert, 0, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;
uniform float mask;

in vec2 fragTexCoord;

out vec4 outputColor;

vec4 toGrayscale(in vec4 color)
{
  float average = (color.r + color.g + color.b) / 3.0;
  return vec4(average, average, average, 1.0);
}

vec4 colorize(in vec4 grayscale, in vec4 color)
{
	return (grayscale * color);
}

void main() {
	vec4 c = vec4(0.2, 0.2, 0.4, 1.0);
	vec4 color = texture(tex, fragTexCoord);
  vec4 grayscale = toGrayscale(color);
	outputColor = mix(color, colorize(grayscale, c), mask);
}
` + "\x00"

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
