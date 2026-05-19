package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const geminiURLTemplate = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"

type GeminiProvider struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewGemini(model, apiKey string) *GeminiProvider {
	return &GeminiProvider{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *GeminiProvider) Translate(ctx context.Context, req TranslateRequest) (TranslateResponse, error) {
	url := fmt.Sprintf(geminiURLTemplate, p.Model, p.APIKey)

	body, err := json.Marshal(map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{{"text": BuildSystemPrompt(req.From, req.To)}},
		},
		"contents": []map[string]any{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": WrapText(req.Text)}},
			},
		},
	})
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("gemini marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return TranslateResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("gemini read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return TranslateResponse{}, fmt.Errorf("gemini %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return TranslateResponse{}, fmt.Errorf("gemini response: %w", err)
	}
	if len(result.Candidates) == 0 {
		return TranslateResponse{}, fmt.Errorf("gemini: empty candidates")
	}
	if len(result.Candidates[0].Content.Parts) == 0 {
		return TranslateResponse{}, fmt.Errorf("gemini: empty parts (finish reason: %s)", result.Candidates[0].FinishReason)
	}
	return TranslateResponse{Text: result.Candidates[0].Content.Parts[0].Text}, nil
}
