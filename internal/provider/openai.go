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

const openAIURL = "https://api.openai.com/v1/chat/completions"

type OpenAIProvider struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewOpenAI(model, apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *OpenAIProvider) Translate(ctx context.Context, req TranslateRequest) (TranslateResponse, error) {
	body, err := json.Marshal(map[string]any{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "system", "content": BuildSystemPrompt(req.From, req.To)},
			{"role": "user", "content": req.Text},
		},
	})
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("openai marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewReader(body))
	if err != nil {
		return TranslateResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TranslateResponse{}, fmt.Errorf("openai read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return TranslateResponse{}, fmt.Errorf("openai %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return TranslateResponse{}, fmt.Errorf("openai response: %w", err)
	}
	if len(result.Choices) == 0 {
		return TranslateResponse{}, fmt.Errorf("openai: empty response")
	}
	return TranslateResponse{Text: result.Choices[0].Message.Content}, nil
}
