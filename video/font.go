package video

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type glyph struct {
	x, y     int
	width    int // Glyph width
	height   int // Glyph height
	advance  int // Glyph advance
	bearingH int // Glyph bearing horizontal
	bearingV int // Glyph bearing vertical
}

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// Known directions.
const (
	LeftToRight Direction = iota // E.g.: Latin
	RightToLeft                  // E.g.: Arabic
	TopToBottom                  // E.g.: Chinese
)

const fontRenderScale = 1

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	glyphs      map[rune]*glyph
	vao         uint32
	vbo         uint32
	program     uint32
	textureID   uint32
	color       Color
	atlasWidth  float32
	atlasHeight float32
	atlas       *image.RGBA
	face        font.Face
	ttf         *truetype.Font
	scale       int32
	lineHeight  float32
	margin      int
	nextX       int
	nextY       int
}

type point [4]float32

type glyphMetrics struct {
	bounds  fixed.Rectangle26_6
	advance fixed.Int26_6
	width   int
	height  int
	ascent  int
	descent int
}

func ceilFixed26_6(v fixed.Int26_6) int32 {
	return int32((v + 63) >> 6)
}

func appendRange(runes []rune, start, end rune) []rune {
	for r := start; r <= end; r++ {
		runes = append(runes, r)
	}
	return runes
}

func defaultRunes() []rune {
	runes := make([]rune, 0, 512)
	runes = appendRange(runes, 32, 126)
	runes = append(runes,
		'◀', '▶', '【', '】',
		'「', '」', '『', '』',
		'。', '、', '・', 'ー',
		'！', '？',
		'（', '）', '［', '］', '｛', '｝',
	)
	runes = appendRange(runes, 0xFF01, 0xFF0F)
	runes = appendRange(runes, 0xFF10, 0xFF19)
	runes = appendRange(runes, 0xFF1A, 0xFF20)
	runes = appendRange(runes, 0xFF21, 0xFF3A)
	runes = appendRange(runes, 0xFF3B, 0xFF40)
	runes = appendRange(runes, 0xFF41, 0xFF5A)
	runes = appendRange(runes, 0xFF5B, 0xFF65)
	runes = appendRange(runes, 0x00A0, 0x017F)
	runes = appendRange(runes, 0x3000, 0x303F)
	runes = appendRange(runes, 0x3040, 0x30FF)
	return runes
}

func getGlyphMetrics(face font.Face, ttf *truetype.Font, ch rune, scale int32) (glyphMetrics, error) {
	gBnd, gAdv, ok := face.GlyphBounds(ch)
	if !ok {
		return glyphMetrics{}, fmt.Errorf("ttf face glyphBounds error")
	}

	gh := ceilFixed26_6(gBnd.Max.Y - gBnd.Min.Y)
	gw := ceilFixed26_6(gBnd.Max.X - gBnd.Min.X)

	if gw == 0 || gh == 0 {
		gBnd = ttf.Bounds(fixed.Int26_6(scale))
		gw = ceilFixed26_6(gBnd.Max.X - gBnd.Min.X)
		gh = ceilFixed26_6(gBnd.Max.Y - gBnd.Min.Y)
		if gw == 0 || gh == 0 {
			gw = 1
			gh = 1
		}
	}

	return glyphMetrics{
		bounds:  gBnd,
		advance: gAdv,
		width:   int(gw),
		height:  int(gh),
		ascent:  int(ceilFixed26_6(-gBnd.Min.Y)),
		descent: int(ceilFixed26_6(gBnd.Max.Y)),
	}, nil
}

func wrapGlyph(x, y *int, glyphWidth, lineHeight, atlasWidth, margin int) {
	if *x+glyphWidth+margin > atlasWidth {
		*x = margin
		*y += lineHeight + margin
	}
}

func appendGlyph(face font.Face, ttf *truetype.Font, ch rune, x, y *int, lineHeight, atlasWidth float32, atlas *image.RGBA, fg *image.Uniform, scale int32, margin int) (*glyph, error) {
	char := new(glyph)

	metrics, err := getGlyphMetrics(face, ttf, ch, scale)
	if err != nil {
		return nil, err
	}

	wrapGlyph(x, y, metrics.width, int(lineHeight), int(atlasWidth), margin)
	if *y+metrics.height+margin > atlas.Rect.Dy() {
		return nil, fmt.Errorf("glyph atlas full")
	}

	char.x = *x
	char.y = *y
	char.width = metrics.width
	char.height = metrics.height
	char.advance = int(metrics.advance)
	char.bearingV = metrics.descent
	char.bearingH = int(metrics.bounds.Min.X) >> 6

	clip := image.Rect(*x, *y, *x+metrics.width, *y+metrics.height)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(float64(scale))
	c.SetClip(clip)
	c.SetDst(atlas)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	px := 0 - (int(metrics.bounds.Min.X) >> 6) + *x
	py := metrics.ascent + *y
	pt := freetype.Pt(px, py)

	_, err = c.DrawString(string(ch), pt)
	*x += metrics.width + margin
	return char, err
}

func getMaxTextureSize() int32 {
	var maxTextureSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	if maxTextureSize == 0 {
		maxTextureSize = 4096
	}
	return maxTextureSize
}

func calculateAtlasSize(face font.Face, ttf *truetype.Font, runes []rune, scale int32, margin int) (float32, float32, float32, error) {
	maxTextureSize := getMaxTextureSize()
	lineHeight := float32(face.Metrics().Height.Ceil())
	if lineHeight < 1 {
		lineHeight = 1
	}

	width := 2048
	if width > int(maxTextureSize) {
		width = int(maxTextureSize)
	}

	for {
		x := margin
		y := margin
		maxY := 0

		for _, r := range runes {
			metrics, err := getGlyphMetrics(face, ttf, r, scale)
			if err != nil {
				continue
			}

			wrapGlyph(&x, &y, metrics.width, int(lineHeight), width, margin)
			x += metrics.width + margin

			if y+int(lineHeight) > maxY {
				maxY = y + int(lineHeight)
			}
		}

		height := maxY + margin
		if height <= int(maxTextureSize) || width >= int(maxTextureSize) {
			if height > int(maxTextureSize) {
				return 0, 0, 0, fmt.Errorf("glyph atlas %dx%d exceeds max texture size %d", width, height, maxTextureSize)
			}
			return float32(width), float32(height), lineHeight, nil
		}

		width *= 2
		if width > int(maxTextureSize) {
			width = int(maxTextureSize)
		}
	}
}

func reserveAtlasHeight(baseHeight int, lineHeight float32, margin int, maxTextureSize int32) int {
	height := baseHeight * 2
	extraRowsHeight := baseHeight + 32*(int(lineHeight)+margin)
	if extraRowsHeight > height {
		height = extraRowsHeight
	}
	if height > int(maxTextureSize) {
		height = int(maxTextureSize)
	}
	return height
}

func (f *Font) uploadGlyphRect(g *glyph) {
	if f.atlas == nil || g == nil {
		return
	}

	offset := g.y*f.atlas.Stride + g.x*4
	if offset < 0 || offset >= len(f.atlas.Pix) {
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(f.atlas.Stride/4))
	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0,
		int32(g.x),
		int32(g.y),
		int32(g.width),
		int32(g.height),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(f.atlas.Pix[offset:]),
	)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// LoadTrueTypeFont builds a set of textures based on a ttf file's glyphs.
func LoadTrueTypeFont(program uint32, r io.Reader, scale int32, dir Direction) (*Font, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	f := &Font{
		glyphs:  make(map[rune]*glyph),
		program: program,
		ttf:     ttf,
		scale:   scale,
	}
	f.SetColor(Color{R: 1, G: 1, B: 1, A: 1})

	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(scale),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	f.face = face

	margin := 4
	runes := defaultRunes()
	f.atlasWidth, f.atlasHeight, f.lineHeight, err = calculateAtlasSize(face, ttf, runes, scale, margin)
	if err != nil {
		return nil, err
	}
	f.atlasHeight = float32(reserveAtlasHeight(int(f.atlasHeight), f.lineHeight, margin, getMaxTextureSize()))
	f.margin = margin

	rect := image.Rect(0, 0, int(f.atlasWidth), int(f.atlasHeight))
	atlas := image.NewRGBA(rect)
	draw.Draw(atlas, atlas.Bounds(), image.Black, image.Point{}, draw.Src)
	f.atlas = atlas

	x := margin
	y := margin
	for _, r := range runes {
		glyph, err := appendGlyph(face, ttf, r, &x, &y, f.lineHeight, f.atlasWidth, atlas, image.White, scale, margin)
		if err != nil {
			continue
		}
		f.glyphs[r] = glyph
	}
	f.nextX = x
	f.nextY = y

	gl.GenTextures(1, &f.textureID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(atlas.Rect.Dx()), int32(atlas.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(atlas.Pix))
	gl.BindTexture(gl.TEXTURE_2D, 0)

	genVertexArrays(1, &f.vao)
	bindVertexArray(f.vao)

	gl.GenBuffers(1, &f.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

	vertAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	bindVertexArray(0)

	gl.UseProgram(f.program)
	gl.Uniform1i(gl.GetUniformLocation(f.program, gl.Str("Texture\x00")), 0)
	gl.UseProgram(0)

	return f, nil
}

func (f *Font) ensureGlyph(r rune) *glyph {
	if glyph, ok := f.glyphs[r]; ok {
		return glyph
	}
	if f.atlas == nil || f.face == nil || f.ttf == nil {
		return nil
	}

	glyph, err := appendGlyph(f.face, f.ttf, r, &f.nextX, &f.nextY, f.lineHeight, f.atlasWidth, f.atlas, image.White, f.scale, f.margin)
	if err != nil {
		log.Printf("font: failed to append U+%04X %q: %v", r, string(r), err)
		return nil
	}
	f.glyphs[r] = glyph
	f.uploadGlyphRect(glyph)
	return glyph
}

// LoadFont loads the specified font at the given scale.
func LoadFont(file string, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	program, err := newProgram(fontVertexShader, fontFragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)
	gl.Uniform2f(gl.GetUniformLocation(program, gl.Str("resolution\x00")), float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)

	atlasScale := int32(float32(scale) / fontRenderScale)
	if atlasScale < 1 {
		atlasScale = 1
	}

	return LoadTrueTypeFont(program, fd, atlasScale, LeftToRight)
}

// SetColor allows you to set the text color to be used when you draw the text.
func (f *Font) SetColor(color Color) {
	f.color = color
}

// UpdateResolution passes the new framebuffer size to the font shader.
func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	gl.UseProgram(f.program)
	gl.Uniform2f(gl.GetUniformLocation(f.program, gl.Str("resolution\x00")), float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)
}

// Print draws a string to the screen.
func (f *Font) Print(x, y float32, scale float32, text string) error {
	indices := []rune(text)
	scale *= fontRenderScale
	x = float32(int(x + 0.5))
	y = float32(int(y + 0.5))

	if len(indices) == 0 {
		return nil
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.UseProgram(f.program)
	gl.Uniform4f(gl.GetUniformLocation(f.program, gl.Str("color\x00")), f.color.R, f.color.G, f.color.B, f.color.A)

	var coords []point
	for i := range indices {
		runeIndex := indices[i]

		ch := f.ensureGlyph(runeIndex)
		if ch == nil {
			ch = f.glyphs['?']
		}
		if ch == nil {
			continue
		}

		xpos := x + float32(ch.bearingH)*scale
		ypos := y - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width) * scale
		h := float32(ch.height) * scale

		x1 := xpos
		x2 := xpos + w
		y1 := ypos
		y2 := ypos + h

		coords = append(coords, point{x1, y1, float32(ch.x) / f.atlasWidth, float32(ch.y) / f.atlasHeight})
		coords = append(coords, point{x2, y1, float32(ch.x+ch.width) / f.atlasWidth, float32(ch.y) / f.atlasHeight})
		coords = append(coords, point{x1, y2, float32(ch.x) / f.atlasWidth, float32(ch.y+ch.height) / f.atlasHeight})
		coords = append(coords, point{x2, y1, float32(ch.x+ch.width) / f.atlasWidth, float32(ch.y) / f.atlasHeight})
		coords = append(coords, point{x1, y2, float32(ch.x) / f.atlasWidth, float32(ch.y+ch.height) / f.atlasHeight})
		coords = append(coords, point{x2, y2, float32(ch.x+ch.width) / f.atlasWidth, float32(ch.y+ch.height) / f.atlasHeight})

		x += float32(ch.advance>>6) * scale
	}

	bindVertexArray(f.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(coords)*16, gl.Ptr(coords), gl.DYNAMIC_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(coords)))
	bindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)

	return nil
}

// Printf formats then draws a string to the screen.
func (f *Font) Printf(x, y float32, scale float32, format string, argv ...interface{}) error {
	return f.Print(x, y, scale, fmt.Sprintf(format, argv...))
}

// Width returns the width of a piece of text in pixels.
func (f *Font) Width(scale float32, text string) float32 {
	var width float32
	scale *= fontRenderScale

	indices := []rune(text)
	if len(indices) == 0 {
		return 0
	}

	for i := range indices {
		runeIndex := indices[i]
		if f.face != nil {
			if advance, ok := f.face.GlyphAdvance(runeIndex); ok {
				width += float32(advance.Round()) * scale
				continue
			}
		}

		if ch, ok := f.glyphs['?']; ok {
			width += float32(ch.advance>>6) * scale
		}
	}

	return width
}
