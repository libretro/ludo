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
	glyphs      map[rune]*glyph
	vao         uint32
	vbo         uint32
	program     uint32
	textureID   uint32 // Holds the glyph texture id.
	color       Color
	atlasWidth  float32
	atlasHeight float32
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

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func appendRange(runes []rune, start, end rune) []rune {
	for r := start; r <= end; r++ {
		runes = append(runes, r)
	}
	return runes
}

func defaultRunes() []rune {
	runes := make([]rune, 0, 22000)
	runes = appendRange(runes, 32, 126)
	runes = append(runes, '◀', '▶', '【', '】')
	runes = appendRange(runes, 0x00A0, 0x017F)
	runes = appendRange(runes, 0x3040, 0x30FF)
	runes = appendRange(runes, 0x4E00, 0x9FAF)
	return runes
}

func getGlyphMetrics(face font.Face, ttf *truetype.Font, ch rune, scale int32) (glyphMetrics, error) {
	gBnd, gAdv, ok := face.GlyphBounds(ch)
	if !ok {
		return glyphMetrics{}, fmt.Errorf("ttf face glyphBounds error")
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

	return glyphMetrics{
		bounds:  gBnd,
		advance: gAdv,
		width:   int(gw),
		height:  int(gh),
		ascent:  gAscent,
		descent: gdescent,
	}, nil
}

func increaseLineHeight(face font.Face, ttf *truetype.Font, ch rune, lineHeight float32, scale int32) float32 {
	metrics, err := getGlyphMetrics(face, ttf, ch, scale)
	if err != nil {
		fmt.Println("ttf face glyphBounds error")
		return lineHeight
	}

	return max(lineHeight, float32(metrics.height))
}

func appendGlyph(face font.Face, ttf *truetype.Font, ch rune, x, y *int, lineHeight, atlasWidth float32, atlas *image.RGBA, fg *image.Uniform, scale int32, margin int) (*glyph, error) {
	char := new(glyph)

	metrics, err := getGlyphMetrics(face, ttf, ch, scale)
	if err != nil {
		return nil, err
	}

	// Set w,h and adv, bearing V and bearing H in char
	char.x = *x
	char.y = *y
	char.width = metrics.width
	char.height = metrics.height
	char.advance = int(metrics.advance)
	char.bearingV = metrics.descent
	char.bearingH = (int(metrics.bounds.Min.X) >> 6)

	clip := image.Rect(*x, *y, *x+metrics.width, *y+metrics.height)

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
	px := 0 - (int(metrics.bounds.Min.X) >> 6) + *x
	py := (metrics.ascent) + *y
	pt := freetype.Pt(px, py)

	*x += metrics.width + margin
	if *x+metrics.width+margin > int(atlasWidth) {
		*x = 0
		*y += int(lineHeight) + margin
	}

	// Draw the text from mask to image
	_, err = c.DrawString(string(ch), pt)
	return char, err
}

func calculateAtlasSize(face font.Face, ttf *truetype.Font, runes []rune, scale int32, margin int) (float32, float32, float32, error) {
	var maxTextureSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	if maxTextureSize == 0 {
		maxTextureSize = 4096
	}

	var lineHeight float32
	for _, r := range runes {
		lineHeight = increaseLineHeight(face, ttf, r, lineHeight, scale)
	}

	width := int(maxTextureSize)
	if width > 4096 {
		width = 4096
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

			if x+metrics.width+margin > width {
				x = margin
				y += int(lineHeight) + margin
			}

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

	margin := 2
	runes := defaultRunes()
	var lineHeight float32

	f.atlasWidth, f.atlasHeight, lineHeight, err = calculateAtlasSize(face, ttf, runes, scale, margin)
	if err != nil {
		return nil, err
	}

	// Create image to draw glyph
	fg, bg := image.White, image.Black
	rect := image.Rect(0, 0, int(f.atlasWidth), int(f.atlasHeight))
	atlas := image.NewRGBA(rect)
	draw.Draw(atlas, atlas.Bounds(), bg, image.Point{}, draw.Src)

	x := margin
	y := margin

	for _, r := range runes {
		glyph, err := appendGlyph(face, ttf, r, &x, &y, lineHeight, f.atlasWidth, atlas, fg, scale, margin)
		if err != nil {
			fmt.Printf("error appending glyph %c: %v\n", r, err)
			continue
		}
		f.glyphs[r] = glyph
	}

	// Generate texture
	gl.GenTextures(1, &f.textureID)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

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

// Print draws a string to the screen.
func (f *Font) Print(x, y float32, scale float32, text string) error {
	indices := []rune(text)

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

// Printf formats then draws a string to the screen.
func (f *Font) Printf(x, y float32, scale float32, format string, argv ...interface{}) error {
	return f.Print(x, y, scale, fmt.Sprintf(format, argv...))
}

// Width returns the width of a piece of text in pixels.
func (f *Font) Width(scale float32, text string) float32 {
	var width float32

	indices := []rune(text)

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
