package video

import (
	"fmt"
	"image"
	"image/draw"
	"io"
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

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	glyphs    map[rune]*glyph
	vao         uint32
	vbo         uint32
	program     uint32
	textureID   uint32 // Holds the glyph texture id.
	color       Color
	atlasWidth  float32
	atlasHeight float32
}

type point [4]float32

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func increaseLineHeight(face font.Face, ch rune, lineHeight float32) float32 {
	gBnd, _, ok := face.GlyphBounds(ch)
	if !ok {
		fmt.Println("ttf face glyphBounds error")
	}
	gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
	lineHeight = max(lineHeight, float32(gh))
	return lineHeight
}

func appendGlyph(face font.Face, ttf *truetype.Font, ch rune, x, y *int, lineHeight, atlasWidth float32, atlas *image.RGBA, fg *image.Uniform, scale int32, margin int) (*glyph, error) {
	char := new(glyph)

	gBnd, gAdv, ok := face.GlyphBounds(ch)
	if !ok {
		return nil, fmt.Errorf("ttf face glyphBounds error")
	}

	gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
	gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

	// If gylph has no dimensions set to a max value
	if gw == 0 || gh == 0 {
		gBnd = ttf.Bounds(fixed.Int26_6(scale))
		gw = int32((gBnd.Max.X - gBnd.Min.X) >> 6)
		gh = int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

		// Above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
		if gw == 0 || gh == 0 {
			gw = 1
			gh = 1
		}
	}

	// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
	gAscent := int(-gBnd.Min.Y) >> 6
	gdescent := int(gBnd.Max.Y) >> 6

	// Set w,h and adv, bearing V and bearing H in char
	char.x = *x
	char.y = *y
	char.width = int(gw)
	char.height = int(gh)
	char.advance = int(gAdv)
	char.bearingV = gdescent
	char.bearingH = (int(gBnd.Min.X) >> 6)

	clip := image.Rect(*x, *y, *x+int(gw), *y+int(gh))

	// Create a freetype context for drawing
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(float64(scale))
	c.SetClip(clip)
	c.SetDst(atlas)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	// Set the glyph dot
	px := 0 - (int(gBnd.Min.X) >> 6) + *x
	py := (gAscent) + *y
	pt := freetype.Pt(px, py)

	*x += int(gw) + margin
	if *x+int(gw)+margin > int(atlasWidth) {
		*x = 0
		*y += int(lineHeight) + margin
	}

	// Draw the text from mask to image
	_, err := c.DrawString(string(ch), pt)
	return char, err
}

// LoadTrueTypeFont builds a set of textures based on a ttf files gylphs
func LoadTrueTypeFont(program uint32, r io.Reader, scale int32, dir Direction) (*Font, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Make Font stuct type
	f := new(Font)
	f.glyphs = make(map[rune]*glyph)
	f.program = program                       // Set shader program
	f.SetColor(Color{R: 1, G: 1, B: 1, A: 1}) // Set default white

	// Create new face
	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(scale),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	var lineHeight float32
	f.atlasWidth = 4096
	f.atlasHeight = 4096
	for r := rune(32); r <= 126; r++ {
		lineHeight = increaseLineHeight(face, r, lineHeight)
	}
	for _, r := range []rune{'◀', '▶', '【', '】'} {
		lineHeight = increaseLineHeight(face, r, lineHeight)
	}
	for r := rune(0x00A0); r <= 0x017F; r++ {
		lineHeight = increaseLineHeight(face, r, lineHeight)
	}
	for r := rune(0x3040); r <= 0x30FF; r++ {
		lineHeight = increaseLineHeight(face, r, lineHeight)
	}
	for r := rune(0x4E00); r <= 0x9FAF; r++ {
		lineHeight = increaseLineHeight(face, r, lineHeight)
	}

	// Create image to draw glyph
	fg, bg := image.White, image.Black
	rect := image.Rect(0, 0, int(f.atlasWidth), int(f.atlasHeight))
	atlas := image.NewRGBA(rect)
	draw.Draw(atlas, atlas.Bounds(), bg, image.Point{}, draw.Src)

	margin := 2
	x := margin
	y := margin

	// Latin
	for r := rune(32); r <= 126; r++ {
		f.glyphs[r], err = appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
		}
	}
	// Some symbols
	for _, r := range []rune{'◀', '▶', '【', '】'} {
		f.glyphs[r], err = appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
		}
	}
	// Extended Latin
	for r := rune(0x00A0); r <= 0x017F; r++ {
		f.glyphs[r], err = appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
		}
	}
	// Japanese Hiragana and Katakana
	for r := rune(0x3040); r <= 0x30FF; r++ {
		f.glyphs[r], err = appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
		}
	}
	// Japanese Kanji
	for r := rune(0x4E00); r <= 0x9FAF; r++ {
		f.glyphs[r], err = appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
		}
	}

	// Generate texture
	gl.GenTextures(1, &f.textureID)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(atlas.Rect.Dx()), int32(atlas.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(atlas.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Configure VAO/VBO for texture quads
	genVertexArrays(1, &f.vao)
	gl.GenBuffers(1, &f.vbo)
	bindVertexArray(f.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

	vertAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	bindVertexArray(0)

	return f, nil
}

// LoadFont loads the specified font at the given scale.
func LoadFont(file string, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// Configure the default font vertex and fragment shaders
	program, err := newProgram(fontVertexShader, fontFragmentShader)
	if err != nil {
		panic(err)
	}

	// Activate corresponding render state
	gl.UseProgram(program)

	// Set screen resolution
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))

	return LoadTrueTypeFont(program, fd, scale, LeftToRight)
}

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(color Color) {
	f.color = color
}

// UpdateResolution passes the new framebuffer size to the font shader
func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	gl.UseProgram(f.program)
	resUniform := gl.GetUniformLocation(f.program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)
}

// Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Printf(x, y float32, scale float32, fs string, argv ...interface{}) error {
	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	// Setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Activate corresponding render state
	gl.UseProgram(f.program)
	// Set text color
	gl.Uniform4f(gl.GetUniformLocation(f.program, gl.Str("textColor\x00")), f.color.R, f.color.G, f.color.B, f.color.A)

	var coords []point

	// Iterate through all glyphs in string
	for i := range indices {
		// Get rune
		runeIndex := indices[i]

		var ch *glyph
		// Skip runes that are not in font glyph map
		if glyph, ok := f.glyphs[runeIndex]; ok {
			ch = glyph
		} else {
			ch = f.glyphs['?'] // fallback
		}

		// Calculate position and size for current rune
		xpos := x - 1 + float32(ch.bearingH)*scale
		ypos := y - 2 - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width+2) * scale
		h := float32(ch.height+2) * scale

		// Set quad positions
		var x1 = xpos
		var x2 = xpos + w
		var y1 = ypos
		var y2 = ypos + h

		coords = append(coords, point{x1, y1, float32(ch.x-1) / f.atlasWidth, float32(ch.y-1) / f.atlasHeight})
		coords = append(coords, point{x2, y1, float32(ch.x+ch.width+1) / f.atlasWidth, float32(ch.y-1) / f.atlasHeight})
		coords = append(coords, point{x1, y2, float32(ch.x-1) / f.atlasWidth, float32(ch.y+ch.height+1) / f.atlasHeight})
		coords = append(coords, point{x2, y1, float32(ch.x+ch.width+1) / f.atlasWidth, float32(ch.y-1) / f.atlasHeight})
		coords = append(coords, point{x1, y2, float32(ch.x-1) / f.atlasWidth, float32(ch.y+ch.height+1) / f.atlasHeight})
		coords = append(coords, point{x2, y2, float32(ch.x+ch.width+1) / f.atlasWidth, float32(ch.y+ch.height+1) / f.atlasHeight})

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
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

// Width returns the width of a piece of text in pixels
func (f *Font) Width(scale float32, fs string, argv ...interface{}) float32 {
	var width float32

	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return 0
	}

	// Iterate through all glyphs in string
	for i := range indices {
		// Get rune
		runeIndex := indices[i]

		// Find rune in glyphs list
		var ch *glyph
		if glyph, ok := f.glyphs[runeIndex]; ok {
			ch = glyph
		} else {
			ch = f.glyphs['?'] // fallback
		}

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	return width
}
