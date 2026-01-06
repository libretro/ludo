package ocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

const (
	openAIURL      = "https://api.openai.com/v1/chat/completions"
	defaultModel   = "gpt-4o-mini"
	defaultPrompt  = "Extract the exact on-screen text. Preserve line breaks. Respond with text only."
	requestTimeout = 30 * time.Second
)

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
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
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

// Preprocess converts the image to a high-contrast black-on-white mask and upscales it.
func Preprocess(img image.Image) image.Image {
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)

	const threshold uint8 = 220
	bin := image.NewGray(gray.Bounds())
	for y := bin.Bounds().Min.Y; y < bin.Bounds().Max.Y; y++ {
		for x := bin.Bounds().Min.X; x < bin.Bounds().Max.X; x++ {
			v := gray.GrayAt(x, y).Y
			if v >= threshold {
				bin.SetGray(x, y, color.Gray{Y: 255})
			} else {
				bin.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	// Upscale to help recognition of small glyphs.
	return imaging.Resize(bin, bin.Bounds().Dx()*4, bin.Bounds().Dy()*4, imaging.NearestNeighbor)
}

// TextFromImage sends the image to OpenAI vision and returns the extracted text.
func TextFromImage(img image.Image) (string, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = defaultModel
	}

	prompt := strings.TrimSpace(os.Getenv("OPENAI_OCR_PROMPT"))
	if prompt == "" {
		prompt = defaultPrompt
	}

	buf := bytes.Buffer{}
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encode image: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	reqBody := chatRequest{
		Model: model,
		Messages: []message{
			{
				Role: "user",
				Content: []contentPart{
					{Type: "text", Text: prompt},
					{
						Type: "image_url",
						ImageURL: &imageURL{
							URL:    "data:image/png;base64," + b64,
							Detail: "high",
						},
					},
				},
			},
		},
		Temperature: 0,
		MaxTokens:   200,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("OpenAI error (%d): %s", resp.StatusCode, string(body))
	}

	var out chatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("OpenAI error: %s (%s)", out.Error.Message, out.Error.Type)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}
