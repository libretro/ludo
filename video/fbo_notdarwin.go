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
	width := video.Geom.MaxWidth
	height := video.Geom.MaxHeight

	log.Printf("[Video]: Initializing HW render (%v x %v).\n", width, height)

	gl.GenFramebuffers(1, &video.fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, video.fboID)

	//gl.GenTextures(1, &video.texID)
	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, int32(width), int32(height))

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, video.texID, 0)

	// Default origin is top left
	video.orthoMat = mgl32.Ortho2D(-1, 1, -1, 1)

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

		if hw.BottomLeftOrigin {
			video.orthoMat = mgl32.Ortho2D(-1, 1, 1, -1)
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

func bindBackbuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
