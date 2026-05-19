package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Languages struct {
		First  string `yaml:"first" mapstructure:"first"`
		Second string `yaml:"second" mapstructure:"second"`
	} `yaml:"languages" mapstructure:"languages"`
	Provider struct {
		Name   string `yaml:"name" mapstructure:"name"`
		Model  string `yaml:"model" mapstructure:"model"`
		APIKey string `yaml:"api_key,omitempty" mapstructure:"api_key"`
	} `yaml:"provider" mapstructure:"provider"`
}

// ConfigPath returns the path to the config file, following the XDG spec (~/.config/tt/config.yaml).
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir: %w", err)
	}
	return filepath.Join(home, ".config", "tt", "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config not found: run tt --init first")
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
