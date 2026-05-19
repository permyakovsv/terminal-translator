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

const anthropicURL = "https://api.anthropic.com/v1/messages"

type AnthropicProvider struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewAnthropic(model, apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *AnthropicProvider) Translate(ctx context.Context, req TranslateRequest) (TranslateResponse, error) {
	body, err := json.Marshal(map[string]any{
		"model":      p.Model,
		"max_tokens": 1024,
		"system":     BuildSystemPrompt(req.From, req.To),
		"messages": []map[string]string{
			{"role": "user", "content": req.Text},
		},
	})
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("anthropic marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicURL, bytes.NewReader(body))
	if err != nil {
		return TranslateResponse{}, err
	}
	httpReq.Header.Set("x-api-key", p.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("anthropic request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("anthropic read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return TranslateResponse{}, fmt.Errorf("anthropic %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return TranslateResponse{}, fmt.Errorf("anthropic response: %w", err)
	}
	if len(result.Content) == 0 || result.Content[0].Type != "text" {
		return TranslateResponse{}, fmt.Errorf("anthropic: unexpected response format")
	}
	return TranslateResponse{Text: result.Content[0].Text}, nil
}
