package main

import (
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/all-core/gl"
)

type color struct {
	r, g, b, a float32
}

func xywhTo4points(x, y, w, h, fbh float32) (x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	x1 = x
	x2 = x
	x3 = x + w
	x4 = x + w
	y1 = fbh - (y + h)
	y2 = fbh - y
	y3 = fbh - (y + h)
	y4 = fbh - y
	return
}

// Draw an image with x, y, w, h
func drawImage(image uint32, x, y, w, h float32, scale float32, c color) {
	_, fbh := window.GetFramebufferSize()
	ffbh := float32(fbh)

	w *= scale
	h *= scale

	x1, y1, x2, y2, x3, y3, x4, y4 := xywhTo4points(x, y, w, h, ffbh)

	drawTextureQuad(image, x1, y1, x2, y2, x3, y3, x4, y4, c)
}

// Draw a colored quad
func drawQuad(x1, y1, x2, y2, x3, y3, x4, y4 float32, c color) {
	drawTextureQuad(video.white, x1, y1, x2, y2, x3, y3, x4, y4, c)
}

// Draw a texture on a polygon
func drawTextureQuad(image uint32, x1, y1, x2, y2, x3, y3, x4, y4 float32, c color) {
	fbw, fbh := window.GetFramebufferSize()
	ffbw := float32(fbw)
	ffbh := float32(fbh)

	var va = []float32{
		//  X, Y, U, V
		x1/ffbw*2 - 1, y1/ffbh*2 - 1, 0, 1, // left-bottom
		x2/ffbw*2 - 1, y2/ffbh*2 - 1, 0, 0, // left-top
		x3/ffbw*2 - 1, y3/ffbh*2 - 1, 1, 1, // right-bottom
		x4/ffbw*2 - 1, y4/ffbh*2 - 1, 1, 0, // right-top
	}

	gl.UseProgram(video.program)
	maskUniform := gl.GetUniformLocation(video.program, gl.Str("mask\x00"))
	gl.Uniform1f(maskUniform, 0)
	gl.Uniform4f(gl.GetUniformLocation(video.program, gl.Str("texColor\x00")), c.r, c.b, c.g, c.a)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.BindVertexArray(video.vao)
	gl.BindTexture(gl.TEXTURE_2D, image)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
}

func textureLoad(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func newImage(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return textureLoad(rgba)
}

func newWhite() uint32 {
	rgba := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.Point{0, 0}, draw.Src)
	return textureLoad(rgba)
}
