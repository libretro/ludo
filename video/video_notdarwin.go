// +build !darwin

package video

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func bindBackbuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func genVertexArrays(n int32, arrays *uint32) {
	gl.GenVertexArrays(n, arrays)
}

func bindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}
