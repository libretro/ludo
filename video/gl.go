package video

import (
	"errors"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

var (
	glslVersion             = "#version 120\n"
	genVertexArrays         = gl.GenVertexArrays
	bindVertexArray         = gl.BindVertexArray
	genFramebuffers         = gl.GenFramebuffers
	bindFramebuffer         = gl.BindFramebuffer
	deleteFramebuffers      = gl.DeleteFramebuffers
	framebufferTexture2D    = gl.FramebufferTexture2D
	checkFramebufferStatus  = gl.CheckFramebufferStatus
	genRenderbuffers        = gl.GenRenderbuffers
	bindRenderbuffer        = gl.BindRenderbuffer
	deleteRenderbuffers     = gl.DeleteRenderbuffers
	renderbufferStorage     = gl.RenderbufferStorage
	framebufferRenderbuffer = gl.FramebufferRenderbuffer
)

// The frontend and hidden core context use the same hints. Select their API
// from the actual context capabilities, including macOS's legacy extensions.
func initGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	window := glfw.GetCurrentContext()
	major, minor := window.GetAttrib(glfw.ContextVersionMajor), window.GetAttrib(glfw.ContextVersionMinor)
	glslVersion = "#version 120\n"
	if major > 3 || (major == 3 && minor >= 2) {
		glslVersion = "#version 150\n"
	} else if major == 3 {
		glslVersion = "#version 130\n"
	}

	genVertexArrays, bindVertexArray = gl.GenVertexArrays, gl.BindVertexArray
	genFramebuffers, bindFramebuffer, deleteFramebuffers = gl.GenFramebuffers, gl.BindFramebuffer, gl.DeleteFramebuffers
	framebufferTexture2D, checkFramebufferStatus = gl.FramebufferTexture2D, gl.CheckFramebufferStatus
	genRenderbuffers, bindRenderbuffer, deleteRenderbuffers = gl.GenRenderbuffers, gl.BindRenderbuffer, gl.DeleteRenderbuffers
	renderbufferStorage, framebufferRenderbuffer = gl.RenderbufferStorage, gl.FramebufferRenderbuffer

	if major < 3 && !glfw.ExtensionSupported("GL_ARB_vertex_array_object") {
		if !glfw.ExtensionSupported("GL_APPLE_vertex_array_object") {
			return errors.New("OpenGL vertex array objects are not supported")
		}
		genVertexArrays, bindVertexArray = gl.GenVertexArraysAPPLE, gl.BindVertexArrayAPPLE
	}
	if major < 3 && !glfw.ExtensionSupported("GL_ARB_framebuffer_object") {
		if !glfw.ExtensionSupported("GL_EXT_framebuffer_object") {
			return errors.New("OpenGL framebuffer objects are not supported")
		}
		genFramebuffers, bindFramebuffer, deleteFramebuffers = gl.GenFramebuffersEXT, gl.BindFramebufferEXT, gl.DeleteFramebuffersEXT
		framebufferTexture2D, checkFramebufferStatus = gl.FramebufferTexture2DEXT, gl.CheckFramebufferStatusEXT
		genRenderbuffers, bindRenderbuffer, deleteRenderbuffers = gl.GenRenderbuffersEXT, gl.BindRenderbufferEXT, gl.DeleteRenderbuffersEXT
		renderbufferStorage, framebufferRenderbuffer = gl.RenderbufferStorageEXT, gl.FramebufferRenderbufferEXT
	}
	return nil
}
