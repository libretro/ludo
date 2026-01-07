package video

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/yomidevs/jmdict-go"
)

const (
	subtitleUserDictPath = "video/userdict.txt"
)

var (
	subtitleTokOnce sync.Once
	subtitleTok     *tokenizer.Tokenizer
	jmdictOnce      sync.Once
	jmdictGlosses   map[string][]senseEntry
	manualGlosses   = map[string][]string{
		"ビッグズ":  {"Biggs"},
		"ビッグス":  {"Biggs"},
		"アバランチ": {"Avalanche"},
		"反神羅":   {"Shinra"},
	}
)

type subtitleToken struct {
	text    string
	color   Color
	reading string
	pron    string
	base    string
	glosses []string
}

type senseEntry struct {
	pos     []string
	glosses []string
}

func ensureSubtitleTokenizer() {
	subtitleTokOnce.Do(func() {
		opts := []tokenizer.Option{}
		if udict, err := dict.NewUserDict(subtitleUserDictPath); err == nil {
			opts = append(opts, tokenizer.UserDict(udict))
		} else {
			fmt.Printf("[subtitle tokenizer] user dict not loaded (%s): %v\n", subtitleUserDictPath, err)
		}

		t, err := tokenizer.New(ipa.Dict(), opts...)
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
		feats := tok.Features()
		pos := ""
		base := surface
		reading := ""
		pron := ""
		if len(feats) > 0 {
			pos = feats[0]
		}
		if len(feats) > 6 {
			base = feats[6]
		}
		if len(feats) > 7 {
			reading = feats[7]
		}
		if len(feats) > 8 {
			pron = feats[8]
		}
		if reading == "" {
			reading = surface
		}
		glosses := lookupGlosses(base, reading, pos)
		res = append(res, subtitleToken{
			text:    surface,
			color:   posColor(pos),
			base:    base,
			reading: reading,
			pron:    pron,
			glosses: glosses,
		})
	}
	return res
}

func lookupGlosses(base, reading, pos string) []string {
	keys := make([]string, 0, 2)
	seen := make(map[string]struct{})
	addKey := func(k string) {
		if k == "" {
			return
		}
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = struct{}{}
		keys = append(keys, k)
	}
	addKey(base)
	addKey(reading)

	for _, k := range keys {
		if g, ok := manualGlosses[k]; ok {
			return g
		}
	}
	jmdictOnce.Do(loadJMDict)
	if jmdictGlosses == nil {
		return nil
	}
	for _, k := range keys {
		if senses, ok := jmdictGlosses[k]; ok {
			if g := bestEnglishGlosses(pos, k, senses); len(g) > 0 {
				return g
			}
		}
	}
	return nil
}

func loadJMDict() {
	path := strings.TrimSpace(os.Getenv("JMDICT_PATH"))
	if path == "" {
		path = "JMdict_e.gz"
	}

	fd, err := os.Open(path)
	if err != nil {
		fmt.Printf("[subtitle tooltip] failed to open JMDict: %v\n", err)
		return
	}
	defer fd.Close()

	var reader = io.Reader(fd)
	if strings.HasSuffix(strings.ToLower(path), ".gz") {
		gr, err := gzip.NewReader(fd)
		if err != nil {
			fmt.Printf("[subtitle tooltip] failed to decompress JMDict: %v\n", err)
			return
		}
		defer gr.Close()
		reader = gr
	}

	dict, _, err := jmdict.LoadJmdict(reader)
	if err != nil {
		fmt.Printf("[subtitle tooltip] failed to load JMDict: %v\n", err)
		return
	}

	glosses := make(map[string][]senseEntry)
	for _, entry := range dict.Entries {
		forms := make([]string, 0, len(entry.Kanji)+len(entry.Readings))
		for _, k := range entry.Kanji {
			if k.Expression != "" {
				forms = append(forms, k.Expression)
			}
		}
		for _, r := range entry.Readings {
			if r.Reading != "" {
				forms = append(forms, r.Reading)
			}
		}

		for _, sense := range entry.Sense {
			eng := englishGlosses(sense)
			if len(eng) == 0 {
				continue
			}
			se := senseEntry{pos: sense.PartsOfSpeech, glosses: eng}
			for _, f := range forms {
				glosses[f] = append(glosses[f], se)
			}
		}
	}

	jmdictGlosses = glosses
	fmt.Printf("[subtitle tooltip] loaded JMDict with %d entries\n", len(jmdictGlosses))
}

func englishGlosses(s jmdict.JmdictSense) []string {
	line := make([]string, 0, len(s.Glossary))
	for _, g := range s.Glossary {
		lang := "eng"
		if g.Language != nil {
			lang = *g.Language
		}
		if lang == "" || lang == "eng" {
			line = append(line, g.Content)
		}
	}
	return line
}

func bestEnglishGlosses(pos, form string, senses []senseEntry) []string {
	var fallback []string
	matched := false
	for _, s := range senses {
		if len(s.glosses) == 0 {
			continue
		}
		if posMatches(pos, s.pos) {
			matched = true
			return s.glosses
		}
		if fallback == nil {
			fallback = s.glosses
		}
	}
	if matched {
		return fallback
	}
	if pos != "" {
		if len(senses) > 0 {
			all := make([]string, 0, len(senses))
			for _, s := range senses {
				all = append(all, s.glosses...)
			}
			fmt.Printf("[subtitle gloss] no POS match for %s (%s); available=%v\n", form, pos, all)
		}
		return nil
	}
	return fallback
}

func posMatches(pos string, sensePOS []string) bool {
	if len(sensePOS) == 0 {
		return false
	}
	switch {
	case strings.HasPrefix(pos, "助詞"):
		for _, p := range sensePOS {
			if strings.Contains(p, "prt") || strings.Contains(p, "particle") {
				return true
			}
		}
	case strings.HasPrefix(pos, "名詞"):
		for _, p := range sensePOS {
			if strings.HasPrefix(p, "n") || strings.Contains(p, "noun") {
				return true
			}
		}
	case strings.HasPrefix(pos, "動詞"):
		for _, p := range sensePOS {
			if strings.HasPrefix(p, "v") || strings.Contains(p, "verb") {
				return true
			}
		}
	case strings.HasPrefix(pos, "形容詞"), strings.HasPrefix(pos, "連体詞"):
		for _, p := range sensePOS {
			if strings.HasPrefix(p, "adj") ||
				strings.Contains(p, "adjective") ||
				strings.Contains(p, "adjectival") ||
				strings.Contains(p, "rentaishi") {
				return true
			}
		}
	case strings.HasPrefix(pos, "副詞"):
		for _, p := range sensePOS {
			if strings.HasPrefix(p, "adv") || strings.Contains(p, "adverb") {
				return true
			}
		}
	default:
		return true
	}
	return false
}

func posColor(pos string) Color {
	switch pos {
	case "名詞": // noun
		return Color{1, 1, 1, 1}
	case "動詞": // verb
		return Color{1.0, 0.75, 0.9, 1} // soft pink
	case "形容詞": // adjective
		return Color{1.0, 0.95, 0.6, 1} // pastel yellow
	case "助詞": // particle
		return Color{0.75, 0.9, 1.0, 1} // light blue
	case "副詞": // adverb
		return Color{0.75, 1.0, 0.8, 1} // soft green
	case "連体詞": // prenominal adjective
		return Color{0.9, 0.95, 0.7, 1}
	case "記号": // symbol/punctuation
		return Color{0.7, 0.7, 0.7, 1}
	default:
		return Color{1, 1, 1, 1}
	}
}

func lineWidth(video *Video, line []subtitleToken, scale float32) float32 {
	var w float32
	for _, tok := range line {
		w += video.Font.Width(scale, "%s", tok.text)
	}
	return w
}

func (video *Video) drawTooltip(tok subtitleToken, x, y, w, ratio, fbw float32) {
	if tok.reading == "" && len(tok.glosses) == 0 {
		return
	}

	lines := []string{}
	if tok.reading != "" {
		lines = append(lines, tok.reading)
	}
	if len(tok.glosses) > 0 {
		lines = append(lines, tok.glosses[0])
	}

	tipScale := 0.4 * ratio
	tipPadding := 12 * ratio
	lineHeight := 38 * ratio
	var tipW float32
	for _, ln := range lines {
		if w := video.Font.Width(tipScale, "%s", ln); w > tipW {
			tipW = w
		}
	}
	tipW += tipPadding * 2
	tipH := float32(len(lines))*lineHeight + tipPadding*2
	tipX := x + (w-tipW)/2
	if tipX < 0 {
		tipX = 0
	}
	if tipX+tipW > fbw {
		tipX = fbw - tipW
	}
	tipY := y - tipH - 8*ratio
	video.DrawRect(tipX, tipY, tipW, tipH, 0.15, Color{0, 0, 0, 1})
	video.Font.SetColor(Color{1, 1, 1, 1})
	for i, ln := range lines {
		lineY := tipY + tipPadding + video.Font.Ascent(tipScale) + lineHeight*float32(i)
		video.Font.Printf(tipX+tipPadding, lineY, tipScale, "%s", ln)
	}
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

// ClearSubtitle hides the currently displayed subtitle.
func (video *Video) ClearSubtitle() {
	video.subtitleText = ""
	video.subtitleLines = nil
	video.subtitleUntil = time.Time{}
}

// SubtitleText returns the current subtitle text.
func (video *Video) SubtitleText() string {
	return video.subtitleText
}

// RenderSubtitle draws the current subtitle, if any.
func (video *Video) RenderSubtitle() {
	if video == nil || video.Font == nil || video.Window == nil {
		return
	}

	if video.subtitleText == "" {
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
	lineHeight := video.Font.LineHeight(scale)
	if lineHeight == 0 {
		return
	}
	ascent := video.Font.Ascent(scale)
	margin := 50 * ratio

	video.Font.UpdateResolution(fbw, fbh)

	cursorX, cursorY := video.Window.GetCursorPos()
	winW, winH := video.Window.GetSize()
	scaleX := float32(fbw) / float32(winW)
	scaleY := float32(fbh) / float32(winH)
	cx := float32(cursorX) * scaleX
	cy := float32(cursorY) * scaleY

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
	rectY := bgY - margin
	topY := rectY + padding

	bgColor := Color{0, 0, 0, 0.75}

	video.DrawRect(bgX, rectY, bgW, bgH, 0.25, bgColor)
	video.Font.SetColor(Color{1, 1, 1, 1})

	var tooltip *subtitleToken
	var tooltipX, tooltipY, tooltipW float32
	for i, line := range lines {
		baseY := topY + ascent + lineHeight*float32(i)
		hitTop := baseY - ascent
		lineWidth := lineWidth(video, line, scale)
		x := float32(fbw)/2 - lineWidth/2
		for _, tok := range line {
			w := video.Font.Width(scale, "%s", tok.text)
			hover := cx >= x && cx <= x+w && cy >= hitTop && cy <= hitTop+lineHeight
			if hover {
				// Draw a subtle highlight behind the token.
				video.DrawRect(x-4*ratio, hitTop-4*ratio, w+8*ratio, lineHeight+8*ratio, 0.15, Color{1, 1, 1, 0.1})
				// fmt.Printf("[subtitle hover] %s\n", tok.text)
				if tok.reading != "" {
					tooltip = &tok
					tooltipX = x
					tooltipY = hitTop
					tooltipW = w
				}
			}

			video.Font.SetColor(tok.color)
			video.Font.Printf(x, baseY, scale, "%s", tok.text)
			x += w
		}
	}

	// Draw tooltip above all subtitle text.
	if tooltip != nil {
		video.drawTooltip(*tooltip, tooltipX, tooltipY, tooltipW, ratio, float32(fbw))
	}
}
