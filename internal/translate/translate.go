package translate

import (
	"context"
	"fmt"

	"github.com/permyakov/tt/internal/config"
	"github.com/permyakov/tt/internal/detect"
	"github.com/permyakov/tt/internal/provider"
)

func Translate(ctx context.Context, cfg *config.Config, text, forceTo string) (string, error) {
	from, to := detect.Direction(text, cfg.Languages.First, cfg.Languages.Second)
	if forceTo != "" {
		to = forceTo
	}

	entry, ok := provider.Find(cfg.Provider.Name)
	if !ok {
		return "", fmt.Errorf("unknown provider %q", cfg.Provider.Name)
	}

	apiKey := entry.ResolveAPIKey(cfg.Provider.APIKey)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider %q: set %s or run tt --init", cfg.Provider.Name, entry.EnvVar)
	}

	p := entry.New(cfg.Provider.Model, apiKey)
	resp, err := p.Translate(ctx, provider.TranslateRequest{Text: text, From: from, To: to})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
