# tt — Terminal Translator

A lightweight, single-binary CLI tool that translates text between two configured languages using LLM providers.

```
$ tt Привіт, світе!
Hello, world!

$ tt Hello, world!
Привіт, світе!
```

## Features

- **Single binary** — no runtime, no dependencies to install
- **Auto-detects direction** — detects source language from Unicode script (Cyrillic, CJK, Hangul, Hiragana, Arabic, Hebrew, Latin)
- **Three LLM providers** — OpenAI, Anthropic, Gemini
- **Any language pair** — configure once, translate both ways instantly
- **Stdin support** — pipe text directly: `echo "hello" | tt`
- **Fast** — direct HTTP to provider APIs, 30s timeout, no overhead

## Installation

### Build from source

```bash
git clone https://github.com/permyakov/terminal-translator
cd terminal-translator
make build
sudo mv tt /usr/local/bin/
```

Requires Go 1.21+.

## Setup

Run the interactive setup once:

```bash
tt --init
```

```
Select first language:  > Ukrainian
Select second language: > English
Select provider:        > OpenAI
Enter API key (leave blank to use OPENAI_API_KEY env var): > sk-...

Configuration saved to ~/.config/tt/config.yaml
```

This writes `~/.config/tt/config.yaml` with your language pair and provider settings.

## API Keys

Keys are read from environment variables first, then from the config file:

| Provider  | Environment variable  |
|-----------|-----------------------|
| OpenAI    | `OPENAI_API_KEY`      |
| Anthropic | `ANTHROPIC_API_KEY`   |
| Gemini    | `GEMINI_API_KEY`      |

```bash
export OPENAI_API_KEY="sk-..."
```

## Usage

Pass text as arguments:

```bash
tt hello world
# → привіт, світе

tt Привіт, світе
# → Hello world
```

Or pipe from stdin:

```bash
echo "hello world" | tt
# → Привіт, світе!

cat article.txt | tt
# → <translated text>
```

## Configuration

`~/.config/tt/config.yaml`:

```yaml
languages:
    first: ru
    second: en
provider:
    name: openai
    model: gpt-4o-mini
    api_key: ""   # optional — prefer env var
```

### Supported providers and default models

| Provider  | Default model              |
|-----------|----------------------------|
| OpenAI    | `gpt-4o-mini`              |
| Anthropic | `claude-haiku-4-5-20251001`|
| Gemini    | `gemini-2.5-flash`         |

You can change `model` in the config file to any model supported by your provider.

## Language direction detection

`tt` detects source language automatically by examining the dominant Unicode script of the input:

| Script detected        | Configured pair example | Direction       |
|------------------------|-------------------------|-----------------|
| Cyrillic               | ru / en                 | ru → en         |
| Latin                  | ru / en                 | en → ru         |
| Hiragana / Katakana    | ja / en                 | ja → en         |
| Han (CJK)              | zh / en                 | zh → en         |
| Hangul                 | ko / en                 | ko → en         |
| Arabic                 | ar / en                 | ar → en         |

For language pairs that share a script (e.g. en / es), the input is assumed to be in the second configured language and translated to the first.

## License

MIT

---

Built with [Claude Code](https://claude.ai/code)
