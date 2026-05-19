package provider

import "os"

// Entry describes a single LLM provider — its metadata and constructor.
type Entry struct {
	Name         string
	DisplayName  string
	DefaultModel string
	EnvVar       string
	New          func(model, apiKey string) Provider
}

// ResolveAPIKey returns the env var value if set, otherwise falls back to configuredKey.
func (e *Entry) ResolveAPIKey(configuredKey string) string {
	if v := os.Getenv(e.EnvVar); v != "" {
		return v
	}
	return configuredKey
}

// Registry is the single source of truth for provider metadata.
var Registry = []Entry{
	{
		Name:         "openai",
		DisplayName:  "OpenAI",
		DefaultModel: "gpt-4o-mini",
		EnvVar:       "OPENAI_API_KEY",
		New:          func(m, k string) Provider { return NewOpenAI(m, k) },
	},
	{
		Name:         "anthropic",
		DisplayName:  "Anthropic",
		DefaultModel: "claude-haiku-4-5-20251001",
		EnvVar:       "ANTHROPIC_API_KEY",
		New:          func(m, k string) Provider { return NewAnthropic(m, k) },
	},
	{
		Name:         "gemini",
		DisplayName:  "Gemini",
		DefaultModel: "gemini-2.5-flash",
		EnvVar:       "GEMINI_API_KEY",
		New:          func(m, k string) Provider { return NewGemini(m, k) },
	},
}

// Find returns the entry for the given internal name (e.g. "openai").
func Find(name string) (*Entry, bool) {
	for i := range Registry {
		if Registry[i].Name == name {
			return &Registry[i], true
		}
	}
	return nil, false
}

// FindByDisplayName returns the entry for the given display name (e.g. "OpenAI").
func FindByDisplayName(displayName string) (*Entry, bool) {
	for i := range Registry {
		if Registry[i].DisplayName == displayName {
			return &Registry[i], true
		}
	}
	return nil, false
}
