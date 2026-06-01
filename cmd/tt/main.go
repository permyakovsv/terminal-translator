package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/permyakov/tt/internal/config"
	"github.com/permyakov/tt/internal/grammar"
	"github.com/permyakov/tt/internal/provider"
	"github.com/permyakov/tt/internal/translate"
)

var initFlag bool
var toFlag string
var grammarFlag bool

var rootCmd = &cobra.Command{
	Use:           "tt [text...]",
	Short:         "Terminal translator — translate text between two configured languages",
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().BoolVar(&initFlag, "init", false, "Run interactive setup")
	rootCmd.Flags().StringVar(&toFlag, "to", "", "Target language code (e.g. uk, en, de); overrides config")
	rootCmd.Flags().BoolVarP(&grammarFlag, "grammar", "g", false, "Check grammar and spelling of source text")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if initFlag {
		return runInit()
	}
	if len(args) > 0 {
		return runTranslate(strings.Join(args, " "), toFlag)
	}
	if isPiped() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		text := strings.TrimSpace(string(data))
		if text == "" {
			return cmd.Help()
		}
		return runTranslate(text, toFlag)
	}
	return cmd.Help()
}

// isPiped reports whether stdin is a pipe or redirect rather than a terminal.
func isPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

var languages = []struct {
	Label string
	Code  string
}{
	{"English", "en"},
	{"Russian", "ru"},
	{"Spanish", "es"},
	{"French", "fr"},
	{"German", "de"},
	{"Chinese", "zh"},
	{"Japanese", "ja"},
	{"Korean", "ko"},
	{"Portuguese", "pt"},
	{"Italian", "it"},
	{"Ukrainian", "uk"},
	{"Polish", "pl"},
}

func langLabels() []string {
	labels := make([]string, len(languages))
	for i, l := range languages {
		labels[i] = l.Label
	}
	return labels
}

func langCode(label string) string {
	for _, l := range languages {
		if l.Label == label {
			return l.Code
		}
	}
	return strings.ToLower(label)
}

func providerDisplayNames() []string {
	names := make([]string, len(provider.Registry))
	for i, e := range provider.Registry {
		names[i] = e.DisplayName
	}
	return names
}

func runInit() error {
	var firstLang, secondLang, providerLabel, apiKey string

	if err := survey.AskOne(&survey.Select{
		Message: "Select first language:",
		Options: langLabels(),
		Default: "Russian",
	}, &firstLang); err != nil {
		return err
	}

	if err := survey.AskOne(&survey.Select{
		Message: "Select second language:",
		Options: langLabels(),
		Default: "English",
	}, &secondLang); err != nil {
		return err
	}

	if err := survey.AskOne(&survey.Select{
		Message: "Select provider:",
		Options: providerDisplayNames(),
		Default: "OpenAI",
	}, &providerLabel); err != nil {
		return err
	}

	entry, _ := provider.FindByDisplayName(providerLabel)

	fmt.Printf("Tip: you can also export %s=%s instead of storing the key.\n", entry.EnvVar, "\"...\"")
	if err := survey.AskOne(&survey.Password{
		Message: fmt.Sprintf("Enter API key (leave blank to use %s env var):", entry.EnvVar),
	}, &apiKey, survey.WithValidator(func(interface{}) error { return nil })); err != nil {
		return err
	}

	cfg := &config.Config{}
	cfg.Languages.First = langCode(firstLang)
	cfg.Languages.Second = langCode(secondLang)
	cfg.Provider.Name = entry.Name
	cfg.Provider.Model = entry.DefaultModel
	cfg.Provider.APIKey = apiKey

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	path, _ := config.ConfigPath()
	fmt.Printf("Configuration saved to %s\n", path)
	return nil
}

func runTranslate(text, forceTo string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	result, err := translate.Translate(context.Background(), cfg, text, forceTo)
	if err != nil {
		return err
	}
	fmt.Println(result)
	if grammarFlag {
		analysis, err := grammar.Check(context.Background(), cfg, text, forceTo)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(analysis)
	}
	return nil
}
