// +build !darwin

package video

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/state"
)

// InitFramebuffer initializes and configures the video frame buffer based on
// information from the HWRenderCallback of the libretro core.
func (video *Video) InitFramebuffer() {
	width := int32(video.Geom.MaxWidth)
	height := int32(video.Geom.MaxHeight)
	video.texWidth = width
	video.texHeight = height

	log.Printf("[Video]: Initializing HW render (%v x %v).\n", width, height)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, video.pixType, video.pixFmt, nil)

	gl.GenFramebuffers(1, &video.fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, video.fboID)
	libretro.SetCurrentFramebufferValue(uintptr(video.fboID))
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, video.texID, 0)

	hw := state.Core.HWRenderCallback
	if hw.Depth {
		gl.GenRenderbuffers(1, &video.rboID)
		gl.BindRenderbuffer(gl.RENDERBUFFER, video.rboID)
		format := gl.DEPTH_COMPONENT24
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

	video.orthoMat = orthoMatrix(hw.BottomLeftOrigin)

	if st := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); st != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalf("[Video] Framebuffer is not complete. Error: %v\n", st)
	}

	bindBackbuffer()
	gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
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
