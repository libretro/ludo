package video

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	openAIChatURL        = "https://api.openai.com/v1/chat/completions"
	defaultGlossModel    = "gpt-4o-mini"
	glossMaxTextLength   = 120
	glossRequestTimeout  = 20 * time.Second
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
		"魔晄":    {"makou"},
		"まこう":   {"makou"},
	}
)

type subtitleToken struct {
	text    string
	color   Color
	reading string
	pron    string
	base    string
	pos     string
	glosses []string
}

type senseEntry struct {
	id      string
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
	tokens := make([]subtitleToken, 0, len(toks))
	candidates := make([][]senseEntry, 0, len(toks))
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
		glosses, senses := lookupGlossCandidates(base, reading)
		tokens = append(tokens, subtitleToken{
			text:    surface,
			color:   posColor(pos),
			base:    base,
			reading: reading,
			pron:    pron,
			pos:     pos,
			glosses: glosses,
		})
		candidates = append(candidates, senses)
	}
	assignGlossesWithContext(tokens, candidates)
	return tokens
}

func lookupGlossCandidates(base, reading string) ([]string, []senseEntry) {
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
			return g, nil
		}
	}
	jmdictOnce.Do(loadJMDict)
	if jmdictGlosses == nil {
		return nil, nil
	}

	result := make([]senseEntry, 0)
	seenIDs := make(map[string]struct{})
	for _, k := range keys {
		if senses, ok := jmdictGlosses[k]; ok {
			for _, s := range senses {
				if s.id != "" {
					if _, exists := seenIDs[s.id]; exists {
						continue
					}
					seenIDs[s.id] = struct{}{}
				}
				result = append(result, s)
			}
		}
	}
	return nil, result
}

func assignGlossesWithContext(tokens []subtitleToken, candidates [][]senseEntry) {
	if len(tokens) == 0 {
		return
	}

	context := make([]contextToken, 0, len(tokens))
	targets := make([]senseSelectionTarget, 0)
	candidateIndex := make(map[int]map[string]int)

	for i, tok := range tokens {
		context = append(context, contextToken{
			Surface: tok.text,
			POS:     tok.pos,
		})

		if len(tok.glosses) > 0 || len(candidates[i]) == 0 {
			continue
		}

		senseCandidates := make([]senseCandidate, 0, len(candidates[i]))
		indexByID := make(map[string]int, len(candidates[i]))
		for j, s := range candidates[i] {
			id := s.id
			if id == "" {
				id = fmt.Sprintf("sense-%d-%d", i, j)
			}
			senseCandidates = append(senseCandidates, senseCandidate{
				ID:    id,
				Gloss: truncateGloss(strings.Join(s.glosses, "; ")),
			})
			indexByID[id] = j
		}

		if len(senseCandidates) == 0 {
			continue
		}

		targets = append(targets, senseSelectionTarget{
			Index:      i,
			Surface:    tok.text,
			Reading:    tok.reading,
			POS:        tok.pos,
			Candidates: senseCandidates,
		})
		candidateIndex[i] = indexByID
	}

	selected := map[int]string{}
	if len(targets) > 0 && strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) != "" {
		ids, err := selectSenseIDsWithLLM(context, targets)
		if err != nil {
			fmt.Printf("[subtitle gloss] OpenAI sense selection failed: %v\n", err)
		} else {
			selected = ids
		}
	}

	for i := range tokens {
		if len(tokens[i].glosses) > 0 || len(candidates[i]) == 0 {
			continue
		}

		if chosenID, ok := selected[i]; ok && chosenID != "" {
			if idxByID, ok := candidateIndex[i]; ok {
				if candIdx, ok := idxByID[chosenID]; ok && candIdx < len(candidates[i]) {
					tokens[i].glosses = candidates[i][candIdx].glosses
					continue
				}
			}
		}

		if g := bestEnglishGlosses(tokens[i].pos, tokens[i].base, candidates[i]); len(g) > 0 {
			tokens[i].glosses = g
			continue
		}
		tokens[i].glosses = candidates[i][0].glosses
	}
}

func selectSenseIDsWithLLM(contextTokens []contextToken, targets []senseSelectionTarget) (map[int]string, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = defaultGlossModel
	}

	prompt := buildSenseSelectionPrompt(contextTokens, targets)
	fmt.Printf("[subtitle gloss] OpenAI request (model=%s):\n%s\n", model, prompt)
	reqBody := chatRequest{
		Model: model,
		Messages: []message{
			{
				Role: "system",
				Content: []contentPart{{
					Type: "text",
					Text: senseSelectorSystemPrompt,
				}},
			},
			{
				Role: "user",
				Content: []contentPart{{
					Type: "text",
					Text: prompt,
				}},
			},
		},
		Temperature: 0,
		MaxTokens:   400,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), glossRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIChatURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("OpenAI error (%d): %s", resp.StatusCode, string(body))
	}

	var out chatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return nil, fmt.Errorf("OpenAI error: %s (%s)", out.Error.Message, out.Error.Type)
	}
	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	content := strings.TrimSpace(out.Choices[0].Message.Content)
	fmt.Printf("[subtitle gloss] OpenAI response:\n%s\n", content)
	selections, err := parseSenseSelections(content)
	if err != nil {
		return nil, fmt.Errorf("decode selection: %w", err)
	}

	result := make(map[int]string, len(selections))
	for _, sel := range selections {
		if sel.ID == nil {
			continue
		}
		result[sel.Index] = *sel.ID
	}

	return result, nil
}

func buildSenseSelectionPrompt(contextTokens []contextToken, targets []senseSelectionTarget) string {
	var b strings.Builder
	b.WriteString("Context tokens in order (index: surface / pos):\n")
	for i, tok := range contextTokens {
		fmt.Fprintf(&b, "%d: %s / %s\n", i, strconv.Quote(tok.Surface), strconv.Quote(tok.POS))
	}

	b.WriteString("\nChoose the best JMDict sense id for each target token using the provided candidates.\n")
	b.WriteString("Respond with JSON array only: [{\"index\":<token index>,\"id\":<candidate id or null>}].\n")
	b.WriteString("Targets:\n")

	for i, t := range targets {
		fmt.Fprintf(&b, "%d. index: %d, surface: %s, reading: %s, pos: %s\n", i+1, t.Index, strconv.Quote(t.Surface), strconv.Quote(t.Reading), strconv.Quote(t.POS))
		b.WriteString("   candidates:\n")
		for _, c := range t.Candidates {
			fmt.Fprintf(&b, "   - id: %s, gloss: %s\n", strconv.Quote(c.ID), strconv.Quote(c.Gloss))
		}
	}

	return b.String()
}

func truncateGloss(gloss string) string {
	gloss = strings.TrimSpace(gloss)
	runes := []rune(gloss)
	if len(runes) <= glossMaxTextLength {
		return gloss
	}
	return string(runes[:glossMaxTextLength]) + "..."
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

		for i, sense := range entry.Sense {
			eng := englishGlosses(sense)
			if len(eng) == 0 {
				continue
			}
			se := senseEntry{
				id:      fmt.Sprintf("%d-%d", entry.Sequence, i+1),
				pos:     sense.PartsOfSpeech,
				glosses: eng,
			}
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

var senseSelectorSystemPrompt = "You are a Japanese sense selector. Choose one JMDict sense ID per target token from the provided candidates using the surrounding context. If no candidate fits, return null. Respond with a JSON array only (no prose, no code fences)."

type contextToken struct {
	Surface string `json:"surface"`
	POS     string `json:"pos,omitempty"`
}

type senseSelectionTarget struct {
	Index      int              `json:"index"`
	Surface    string           `json:"surface"`
	Reading    string           `json:"reading,omitempty"`
	POS        string           `json:"pos,omitempty"`
	Candidates []senseCandidate `json:"candidates"`
}

type senseCandidate struct {
	ID    string `json:"id"`
	Gloss string `json:"gloss"`
}

type senseSelectionResponse struct {
	Index int     `json:"index"`
	ID    *string `json:"id"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string        `json:"role"`
	Content []contentPart `json:"content"`
}

type contentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func parseSenseSelections(raw string) ([]senseSelectionResponse, error) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start >= 0 && end > start {
		raw = raw[start : end+1]
	}
	var selections []senseSelectionResponse
	if err := json.Unmarshal([]byte(raw), &selections); err != nil {
		return nil, err
	}
	return selections, nil
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
