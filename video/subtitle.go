package video

import (
	"strings"
	"sync"
	"time"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

var (
	subtitleTokOnce sync.Once
	subtitleTok     *tokenizer.Tokenizer
)

type subtitleToken struct {
	text  string
	color Color
}

func ensureSubtitleTokenizer() {
	subtitleTokOnce.Do(func() {
		t, err := tokenizer.New(ipa.Dict())
		if err != nil {
			panic(err)
		}
		subtitleTok = t
	})
}

func tokenizeLines(text string) [][]subtitleToken {
	ensureSubtitleTokenizer()
	lines := strings.Split(text, "\n")
	res := make([][]subtitleToken, 0, len(lines))
	for _, line := range lines {
		tokens := tokenizeLine(line)
		if len(tokens) == 0 && strings.TrimSpace(line) != "" {
			tokens = []subtitleToken{{text: line, color: Color{1, 1, 1, 1}}}
		}
		res = append(res, tokens)
	}
	return res
}

func tokenizeLine(line string) []subtitleToken {
	if subtitleTok == nil {
		ensureSubtitleTokenizer()
	}

	toks := subtitleTok.Tokenize(line)
	res := make([]subtitleToken, 0, len(toks))
	for _, tok := range toks {
		if tok.Class == tokenizer.DUMMY {
			continue
		}
		surface := tok.Surface
		if surface == "" {
			continue
		}
		pos := ""
		feats := tok.Features()
		if len(feats) > 0 {
			pos = feats[0]
		}
		res = append(res, subtitleToken{
			text:  surface,
			color: posColor(pos),
		})
	}
	return res
}

func posColor(pos string) Color {
	switch pos {
	case "名詞": // noun
		return Color{1, 1, 1, 1}
	case "動詞": // verb
		return Color{0.6, 0.9, 1, 1}
	case "形容詞": // adjective
		return Color{1, 0.8, 1, 1}
	case "助詞": // particle
		return Color{0.8, 0.8, 0.8, 1}
	case "副詞": // adverb
		return Color{1, 0.95, 0.7, 1}
	case "連体詞": // prenominal adjective
		return Color{0.9, 1, 0.7, 1}
	case "記号": // symbol/punctuation
		return Color{0.7, 0.7, 0.7, 1}
	default:
		return Color{1, 1, 1, 1}
	}
}

func lineWidth(video *Video, line []subtitleToken, scale float32) float32 {
	var w float32
	for _, tok := range line {
		w += video.Font.Width(scale, tok.text)
	}
	return w
}

func maxLineWidth(video *Video, lines [][]subtitleToken, scale float32) float32 {
	var max float32
	for _, line := range lines {
		if lw := lineWidth(video, line, scale); lw > max {
			max = lw
		}
	}
	return max
}

// SetSubtitle displays text at the bottom of the screen for a given duration.
func (video *Video) SetSubtitle(text string, duration time.Duration) {
	video.subtitleText = strings.TrimSpace(text)
	if video.subtitleText == "" {
		video.subtitleUntil = time.Time{}
		video.subtitleLines = nil
		return
	}

	video.subtitleLines = tokenizeLines(video.subtitleText)
	video.subtitleUntil = time.Now().Add(duration)
}

// RenderSubtitle draws the current subtitle, if any.
func (video *Video) RenderSubtitle() {
	if video == nil || video.Font == nil || video.Window == nil {
		return
	}

	if video.subtitleText == "" || time.Now().After(video.subtitleUntil) {
		video.subtitleText = ""
		return
	}

	fbw, fbh := video.Window.GetFramebufferSize()
	if fbw == 0 || fbh == 0 {
		return
	}

	lines := video.subtitleLines
	if len(lines) == 0 {
		lines = tokenizeLines(video.subtitleText)
	}
	ratio := float32(fbw) / 1920
	scale := 0.6 * ratio
	padding := 16 * ratio
	lineHeight := 64 * ratio
	margin := 50 * ratio

	video.Font.UpdateResolution(fbw, fbh)

	previousColor := video.Font.color
	defer video.Font.SetColor(previousColor)

	maxWidth := maxLineWidth(video, lines, scale)

	bgW := maxWidth + padding*2
	minBgW := float32(fbw) * 0.75
	if bgW < minBgW {
		bgW = minBgW
	}
	bgX := (float32(fbw) - bgW) / 2
	bgH := lineHeight*float32(len(lines)) + padding*2
	bgY := float32(fbh) - bgH
	topY := bgY + padding

	bgColor := Color{0, 0, 0, 0.75}

	video.DrawRect(bgX, bgY-margin, bgW, bgH, 0.25, bgColor)
	video.Font.SetColor(Color{1, 1, 1, 1})

	for i, line := range lines {
		y := topY + lineHeight*float32(i)
		lineWidth := lineWidth(video, line, scale)
		x := float32(fbw)/2 - lineWidth/2
		for _, tok := range line {
			video.Font.SetColor(tok.color)
			video.Font.Printf(x, y, scale, tok.text)
			x += video.Font.Width(scale, tok.text)
		}
	}
}
