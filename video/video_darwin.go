// +build darwin

package video

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func bindBackbuffer() {
	gl.BindFramebufferEXT(gl.FRAMEBUFFER_EXT, 0)
}

func genVertexArrays(n int32, arrays *uint32) {
	gl.GenVertexArraysAPPLE(n, arrays)
}

func bindVertexArray(array uint32) {
	gl.BindVertexArrayAPPLE(array)
}
