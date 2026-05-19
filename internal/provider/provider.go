package provider

import (
	"context"
	"fmt"
)

type Provider interface {
	Translate(ctx context.Context, req TranslateRequest) (TranslateResponse, error)
}

type TranslateRequest struct {
	Text string
	From string
	To   string
}

type TranslateResponse struct {
	Text string
}

func BuildSystemPrompt(from, to string) string {
	return fmt.Sprintf(
		"You are a translation engine.\n\nTranslate the input from %s to %s.\n\nReturn only translated text.\n\nPreserve meaning, tone, idioms, punctuation and formatting.",
		from, to,
	)
}
