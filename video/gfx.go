package video

import (
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/all-core/gl"
)

// Color is an RGBA type that we use in the menu
type Color struct {
	R, G, B, A float32
}

// XYWHTo4points converts coordinates from (x, y, width, height) to (x1, y1, x2, y2, x3, y3, x4, y4)
func XYWHTo4points(x, y, w, h, fbh float32) (x1, y1, x2, y2, x3, y3, x4, y4 float32) {
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

// DrawImage draws an image with x, y, w, h
func (video *Video) DrawImage(image uint32, x, y, w, h float32, scale float32, c Color) {
	video.drawTexturedQuad(image, x, y, w, h, scale, c)
}

// DrawRect draws a colored rectangle
func (video *Video) DrawRect(x, y, w, h float32, scale float32, c Color) {
	video.drawTexturedQuad(video.white, x, y, w, h, scale, c)
}

func (video *Video) vertexArray(x, y, w, h, scale float32) []float32 {
	fbw, fbh := video.Window.GetFramebufferSize()
	ffbw := float32(fbw)
	ffbh := float32(fbh)

	w *= scale
	h *= scale

	x1, y1, x2, y2, x3, y3, x4, y4 := XYWHTo4points(x, y, w, h, ffbh)

	return []float32{
		//  X, Y, U, V
		x1/ffbw*2 - 1, y1/ffbh*2 - 1, 0, 1, // left-bottom
		x2/ffbw*2 - 1, y2/ffbh*2 - 1, 0, 0, // left-top
		x3/ffbw*2 - 1, y3/ffbh*2 - 1, 1, 1, // right-bottom
		x4/ffbw*2 - 1, y4/ffbh*2 - 1, 1, 0, // right-top
	}
}

// DrawBorder draws a colored rectangle border
func (video *Video) DrawBorder(x, y, w, h, borderWidth float32, c Color) {

	va := video.vertexArray(x, y, w, h, 1.0)

	gl.UseProgram(video.borderProgram)
	gl.Uniform1f(gl.GetUniformLocation(video.borderProgram, gl.Str("border_width\x00")), borderWidth)
	gl.Uniform4f(gl.GetUniformLocation(video.borderProgram, gl.Str("color\x00")), c.R, c.G, c.B, c.A)
	gl.Uniform2f(gl.GetUniformLocation(video.borderProgram, gl.Str("size\x00")), w, h)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindVertexArray(video.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
}

// Draw a texture on a polygon
func (video *Video) drawTexturedQuad(image uint32, x, y, w, h, scale float32, c Color) {

	va := video.vertexArray(x, y, w, h, scale)

	gl.UseProgram(video.demulProgram)
	maskUniform := gl.GetUniformLocation(video.demulProgram, gl.Str("mask\x00"))
	gl.Uniform1f(maskUniform, 0)
	gl.Uniform4f(gl.GetUniformLocation(video.demulProgram, gl.Str("texColor\x00")), c.R, c.G, c.B, c.A)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindVertexArray(video.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, image)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
}

// DrawRoundedRect draws a rectangle with rounded corners
func (video *Video) DrawRoundedRect(x, y, w, h, r float32, c Color) {

	va := video.vertexArray(x, y, w, h, 1.0)

	gl.UseProgram(video.roundedProgram)
	gl.Uniform4f(gl.GetUniformLocation(video.roundedProgram, gl.Str("color\x00")), c.R, c.G, c.B, c.A)
	gl.Uniform1f(gl.GetUniformLocation(video.roundedProgram, gl.Str("radius\x00")), r)
	gl.Uniform2f(gl.GetUniformLocation(video.roundedProgram, gl.Str("size\x00")), w, h)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindVertexArray(video.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
}

// DrawCircle draws a circle
func (video *Video) DrawCircle(x, y, r float32, c Color) {

	va := video.vertexArray(x-r, y-r, r*2, r*2, 1.0)

	gl.UseProgram(video.circleProgram)
	gl.Uniform4f(gl.GetUniformLocation(video.circleProgram, gl.Str("color\x00")), c.R, c.G, c.B, c.A)
	gl.Uniform1f(gl.GetUniformLocation(video.circleProgram, gl.Str("radius\x00")), r)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindVertexArray(video.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, video.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
}

func textureLoad(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
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
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

// NewImage opens an image file, upload it the the GPU and returns the texture id
func NewImage(file string) uint32 {
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
