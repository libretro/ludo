// +build darwin

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

	gl.BindTexture(gl.TEXTURE_2D, video.texID)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, int32(width), int32(height))

	gl.GenFramebuffersEXT(1, &video.fboID)
	gl.BindFramebufferEXT(gl.FRAMEBUFFER_EXT, video.fboID)
	gl.FramebufferTexture2DEXT(gl.FRAMEBUFFER_EXT, gl.COLOR_ATTACHMENT0_EXT, gl.TEXTURE_2D, video.texID, 0)

	hw := state.Global.Core.HWRenderCallback
	if hw.Depth {
		gl.GenRenderbuffersEXT(1, &video.rboID)
		gl.BindRenderbufferEXT(gl.RENDERBUFFER_EXT, video.rboID)
		format := gl.DEPTH_COMPONENT16
		if hw.Stencil {
			format = gl.DEPTH24_STENCIL8_EXT
		}
		gl.RenderbufferStorageEXT(gl.RENDERBUFFER_EXT, uint32(format), int32(width), int32(height))
		gl.BindRenderbufferEXT(gl.RENDERBUFFER_EXT, 0)

		gl.FramebufferRenderbufferEXT(gl.FRAMEBUFFER_EXT, gl.DEPTH_ATTACHMENT_EXT, gl.RENDERBUFFER_EXT, video.rboID)
		if hw.Stencil {
			gl.FramebufferRenderbufferEXT(gl.FRAMEBUFFER_EXT, gl.STENCIL_ATTACHMENT_EXT, gl.RENDERBUFFER_EXT, video.rboID)
		}
	}

	// Default origin is top left
	video.orthoMat = mgl32.Ortho2D(-1, 1, -1, 1)
	if hw.BottomLeftOrigin {
		video.orthoMat = mgl32.Ortho2D(-1, 1, 1, -1)
	}

	if st := gl.CheckFramebufferStatusEXT(gl.FRAMEBUFFER_EXT); st != gl.FRAMEBUFFER_COMPLETE_EXT {
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
	gl.BindFramebufferEXT(gl.FRAMEBUFFER_EXT, 0)
}
