# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**tt** (Terminal Translator) is a lightweight, single-binary CLI tool that translates text between two configured languages using LLM providers (OpenAI, Anthropic, Gemini). It auto-detects the source language from Unicode script detection and supports both command-line arguments and stdin piping.

## Build & Development Commands

```bash
make build   # produces ./tt
make dist    # cross-platform binaries in dist/
make clean   # remove build artifacts
```

## Coding Principles

- **Keep it simple**: prefer the straightforward solution over the clever one.
- **Keep it clean**: consistent formatting, clear naming, no dead code.
- **No unnecessary abstractions**: don't generalize until there are at least three concrete cases.
- **No speculative features**: implement what is asked, nothing more.
- **Small functions**: each function does one thing.
- **No comments on obvious code**: only comment non-obvious invariants or workarounds.
- **No SDK dependencies for providers**: all HTTP requests are hand-crafted to keep the binary small and dependency-free.
- **Fail fast**: return errors early; avoid deep nesting.
