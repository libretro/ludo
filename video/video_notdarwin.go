// +build !darwin

package video

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/libretro/ludo/state"
)

// InitFramebuffer initializes and configures the video frame buffer based on
// informations from the HWRenderCallback of the libretro core.
func (video *Video) InitFramebuffer() {
	width := int32(video.Geom.MaxWidth)
	height := int32(video.Geom.MaxHeight)

	log.Printf("[Video]: Initializing HW render (%v x %v).\n", width, height)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, width, height)

	gl.GenFramebuffers(1, &video.fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, video.fboID)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, video.texID, 0)

	hw := state.Global.Core.HWRenderCallback
	if hw.Depth {
		gl.GenRenderbuffers(1, &video.rboID)
		gl.BindRenderbuffer(gl.RENDERBUFFER, video.rboID)
		format := gl.DEPTH_COMPONENT16
		if hw.Stencil {
			format = gl.DEPTH24_STENCIL8
		}
		gl.RenderbufferStorage(gl.RENDERBUFFER, uint32(format), width, height)
		gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

		if hw.Stencil {
			gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, video.rboID)
		} else {
			gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, video.rboID)
		}
	}

	// Default origin is top left
	video.orthoMat = mgl32.Ortho2D(-1, 1, -1, 1)
	if hw.BottomLeftOrigin {
		video.orthoMat = mgl32.Ortho2D(-1, 1, 1, -1)
	}

	if st := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); st != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalf("[Video] Framebuffer is not complete. Error: %v\n", st)
	}

	bindBackbuffer()

	gl.ClearColor(0, 0, 0, 1)
	if hw.Depth && hw.Stencil {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	} else if hw.Depth {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	} else {
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}
}

func bindBackbuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func genVertexArrays(n int32, arrays *uint32) {
	gl.GenVertexArrays(n, arrays)
}

func bindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}
