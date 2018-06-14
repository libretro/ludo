package main

import (
	"fmt"
	"log"

	"strings"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kivutar/glfont"
	"github.com/kivutar/go-playthemall/libretro"
)

var window *glfw.Window

var video struct {
	program uint32
	vao     uint32
	vbo     uint32
	texID   uint32
	white   uint32
	pitch   int32
	pixFmt  uint32
	pixType uint32
	bpp     int32
	geom    libretro.GameGeometry
	font    *glfont.Font
}

func videoSetPixelFormat(format uint32) bool {
	if g.verbose {
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

	return true
}

func updateMaskUniform() {
	maskUniform := gl.GetUniformLocation(video.program, gl.Str("mask\x00"))
	if g.menuActive {
		gl.Uniform1f(maskUniform, 1.0)
	} else {
		gl.Uniform1f(maskUniform, 0.0)
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

func coreRatioViewport(win *glfw.Window, fbWidth int, fbHeight int) {
	// Scale the content to fit in the viewport.
	fbw := float32(fbWidth)
	fbh := float32(fbHeight)
	w := fbw
	h := fbw / float32(video.geom.AspectRatio)
	if h > fbh {
		w = fbh * float32(video.geom.AspectRatio)
		h = fbh
	}

	// Place the content in the middle of the window.
	x := (fbw - w) / 2
	y := (fbh - h) / 2

	x1, y1, x2, y2, x3, y3, x4, y4 := xywhToCoords(x, y, w, h, fbh)

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

func fullscreenViewport() {
	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))
}

func videoConfigure(geom libretro.GameGeometry, fullscreen bool) {
	video.geom = geom

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var width, height int
	var m *glfw.Monitor

	if fullscreen {
		m = glfw.GetMonitors()[settings.VideoMonitorIndex]
		vms := m.GetVideoModes()
		vm := vms[len(vms)-1]
		width = vm.Width
		height = vm.Height
	} else {
		nwidth, nheight := resizeToAspect(geom.AspectRatio, float64(geom.BaseWidth), float64(geom.BaseHeight))
		width = int(nwidth * float64(settings.VideoScale))
		height = int(nheight * float64(settings.VideoScale))
	}

	if window != nil {
		window.Destroy()
	}

	var err error
	window, err = glfw.CreateWindow(width, height, "Play Them All", m, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Force a minimum size for the window.
	window.SetSizeLimits(160, 120, glfw.DontCare, glfw.DontCare)

	if fullscreen {
		window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	fbw, fbh := window.GetFramebufferSize()
	coreRatioViewport(window, fbw, fbh)

	//load font (fontfile, font scale, window width, window height
	video.font, err = glfont.LoadFont("assets/font.ttf", int32(64), fbw, fbh)
	if err != nil {
		panic(err)
	}

	if g.verbose {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("[Video]: OpenGL version:", version)
	}

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

	if video.pixFmt == 0 {
		video.pixFmt = gl.UNSIGNED_SHORT_5_5_5_1
	}

	gl.GenTextures(1, &video.texID)

	gl.ActiveTexture(gl.TEXTURE0)
	if video.texID == 0 && g.verbose {
		log.Println("[Video]: Failed to create the video texture")
	}

	video.pitch = int32(geom.BaseWidth) * video.bpp

	gl.BindTexture(gl.TEXTURE_2D, video.texID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	contextReset()
}

func renderNotifications() {
	fullscreenViewport()
	w, h := window.GetFramebufferSize()
	video.font.UpdateResolution(w, h)
	for i, n := range notifications {
		video.font.SetColor(1.0, 1.0, 0.0, float32(n.frames)/120.0)
		video.font.Printf(80, float32(h-80*len(notifications)+80*i), 0.7, n.message)
	}
}

// Render the current frame
func videoRender() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	fbw, fbh := window.GetFramebufferSize()
	coreRatioViewport(window, fbw, fbh)

	gl.UseProgram(video.program)
	updateMaskUniform()
	gl.Uniform4f(gl.GetUniformLocation(video.program, gl.Str("texColor\x00")), 1, 1, 1, 1)

	gl.BindVertexArray(video.vao)

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// Refresh the texture framebuffer
func videoRefresh(data unsafe.Pointer, width int32, height int32, pitch int32) {
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, video.pixType, video.pixFmt, nil)

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
uniform vec4 texColor;

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
	outputColor = texColor * mix(color, colorize(grayscale, c), mask);
}
` + "\x00"

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
