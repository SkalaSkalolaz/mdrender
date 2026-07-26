# Gogitor

[Русский](README_RU.md)

Gogitor is a terminal AI coding assistant for Go projects. It works both as an interactive TUI and as a classic CLI tool. Gogitor connects to a local or remote LLM (Ollama or OpenAI-compatible API), understands your task, creates or modifies Go files, validates the result in a temporary sandbox using `go build` and `go test`, and can automatically commit changes to Git.

> Current version: `0.9.12`

---

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Main modes](#main-modes)
- [CLI usage](#cli-usage)
- [TUI mode](#tui-mode)
- [Configuration](#configuration)
- [Examples](#examples)
- [How code generation works](#how-code-generation-works)
- [Safety](#safety)
- [Project structure](#project-structure)
- [Troubleshooting](#troubleshooting)

---

## Features

- **Interactive TUI**
  - Bubble Tea based terminal interface.
  - Markdown rendering for chat, analysis, and search answers.
  - Command history, autocomplete for built-in commands, mouse text selection mode.

- **CLI and pipes**
  - Unix-friendly CLI.
  - `--raw` mode for piping generated code directly into files or other tools.
  - `--json` mode for machine-readable results.

- **Intent routing**
  - Gogitor can automatically choose a mode:
    - `code`
    - `analyze`
    - `search`
    - `run`
    - `test`
    - `git`
    - `chat`

- **Code generation and modification**
  - Creates new Go files.
  - Modifies existing files.
  - Supports patch mode with `SEARCH/REPLACE` blocks.
  - Falls back to full-file rewrite when patching is unsafe or fails.

- **Sandbox validation**
  - Copies the project to a temporary directory.
  - Runs:
    - `go mod init` if needed
    - `go mod tidy` when external imports are detected
    - `gofmt`
    - `go build`
    - `go test -v -cover`
  - Applies changes to the real project only after successful validation.

- **Tests and coverage**
  - Parses passed/failed tests.
  - Extracts failure details.
  - Shows coverage information when available.

- **Multi-agent planning**
  - For complex tasks, Gogitor can split the task into subtasks and execute them step by step.

- **Git integration**
  - `status`
  - `diff`
  - `commit`
  - `init`
  - `log`
  - `checkout`
  - `branch`
  - `merge`
  - Optional automatic commit after successful code changes.

- **Web search**
  - DuckDuckGo-based search mode.
  - LLM summarizes search results into an answer.

- **Task files**
  - Run tasks from `.txt` or `.md` files.

- **Security guards**
  - Path traversal protection for file changes.
  - Temporary sandbox for build/test validation.

---

## Requirements

- **Go 1.21+**
- **Git** — optional, but strongly recommended
- **Ollama** or an **OpenAI-compatible API endpoint**
- **Unix-like OS**: Linux, macOS, or WSL
- Network access:
  - for downloading Go dependencies
  - for web search mode
  - for remote LLM APIs

---

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/YOUR-USERNAME/gogitor.git
cd gogitor
go mod tidy
go build -o gogitor .
```

Then run:

```bash
./gogitor --help
```

You can also install it into your `GOPATH` manually if you prefer:

```bash
go build -o gogitor .
sudo mv gogitor /usr/local/bin/
```

---

## Quick start

### Start interactive TUI

```bash
./gogitor
```

Or explicitly specify provider and model:

```bash
./gogitor tui --provider ollama --model gemma3:4b
```

### Generate code from CLI

```bash
./gogitor code "write a Go program that multiplies two matrices"
```

### Ask a question

```bash
./gogitor ask "explain context.Context in Go"
```

### Analyze current project

```bash
./gogitor analyze "find potential bugs and suggest improvements"
```

### Run tests

```bash
./gogitor test
```

### Check Git status

```bash
./gogitor git status
```

---

## Main modes

### `code`

Creates or modifies code.

Example:

```bash
./gogitor code "add a HTTP health endpoint to main.go"
```

### `ask`

Chat mode. Good for general Go questions.

```bash
./gogitor ask "what is the difference between make and new in Go?"
```

### `analyze`

Reads project files and analyzes them without modifying anything.

```bash
./gogitor analyze "review this project and point out architecture issues"
```

### `search`

Web search plus LLM summary.

```bash
./gogitor search "latest Go release notes"
```

### `run`

Runs the Go program in a sandbox.

```bash
./gogitor run
```

Or run a specific file directory:

```bash
./gogitor run cmd/server/main.go
```

### `test`

Runs tests in a sandbox.

```bash
./gogitor test
```

### `git`

Git helper.

```bash
./gogitor git status
./gogitor git diff
./gogitor git commit
./gogitor git log
```

---

## CLI usage

```text
gogitor
gogitor tui [flags]
gogitor code <task> [flags]
gogitor task <path/to/file.txt|file.md> [flags]
gogitor file <path/to/file.txt|file.md> [flags]
gogitor ask <question> [flags]
gogitor analyze <question> [flags]
gogitor search <query> [flags]
gogitor run [file] [flags]
gogitor test [flags]
gogitor git <subcommand> [flags]
gogitor doctor [flags]
gogitor help
gogitor version
```

### Commands

| Command | Description |
|---|---|
| `gogitor` | Starts TUI if stdin is a terminal; if input is piped, treats stdin as an `ask` request |
| `gogitor tui` | Starts interactive TUI |
| `gogitor code <task>` | Direct code generation/modification mode |
| `gogitor task <file>` | Executes a task from a `.txt` or `.md` file |
| `gogitor file <file>` | Alias of `task` |
| `gogitor ask <question>` | Chat mode |
| `gogitor analyze <question>` | Analyze project files |
| `gogitor search <query>` | Web search mode |
| `gogitor run [file]` | Run Go code in sandbox |
| `gogitor test` | Run tests in sandbox |
| `gogitor git <subcommand>` | Git operations |
| `gogitor doctor` | Show configuration and diagnostics |
| `gogitor help` | Show help |
| `gogitor version` | Show version |

### Common flags

| Flag | Short | Description |
|---|---:|---|
| `--provider <name>` | `-p` | LLM provider |
| `--model <model>` | `-m` | Model name |
| `--key <key>` | `-k` | API key |
| `--repo <path>` | `-r` | Project root directory |
| `--debug` |  | Enable debug logging |
| `--raw` |  | Output only result content |
| `--pretty` |  | Force human-readable output |
| `--help` | `-h` | Show help |

### Code mode flags

| Flag | Description |
|---|---|
| `--dry-run` | Validate changes but do not apply them |
| `--no-commit` | Disable automatic Git commit |
| `--no-tests` | Skip tests |
| `--json` | Print result as JSON |

### Task file flags

| Flag | Description |
|---|---|
| `--code` | Force code mode instead of automatic intent detection |
| `--json` | Print result as JSON |

### Git subcommands

| Subcommand | Description |
|---|---|
| `status` | Show Git status |
| `diff` | Show working tree diff |
| `commit` | Commit all changes |
| `init` | Initialize Git repository |
| `log` | Show commit history |
| `checkout <hash-or-branch>` | Checkout commit or branch |
| `branch` | List branches |
| `branch <name>` | Create branch |
| `branch -d <name>` | Delete branch |
| `branch -D <name>` | Force delete branch |
| `merge <branch>` | Merge branch into current branch |

---

## TUI mode

Run:

```bash
./gogitor
```

or:

```bash
./gogitor tui
```

### Built-in commands

Inside TUI you can type:

| Command | Description |
|---|---|
| `:help` | Show help |
| `:clear` | Clear in-memory conversation context |
| `:cls` | Clear TUI screen |
| `:code <task>` | Code mode |
| `:ask <question>` | Chat mode |
| `:analyze <task>` | Analyze mode |
| `:search <query>` | Search mode |
| `:run [file]` | Run project |
| `:test` | Run tests |
| `:git <subcommand>` | Git operations |
| `:quit` or `:q` | Quit |

### Keyboard shortcuts

| Key | Action |
|---|---|
| `Enter` | Send input |
| `Alt+Enter` | New line |
| `Tab` | Switch focus between input and output / autocomplete |
| `Esc` | Return to input from output focus |
| `PgUp` / `PgDn` | History navigation or output scrolling |
| `F2` | Toggle mouse text selection mode |
| `Ctrl+C` | Cancel current task or quit |

---

## Configuration

Gogitor reads configuration in this order:

1. Built-in defaults
2. Global config file: `~/.gogitor/config.json`
3. Environment variables
4. Local project config: `.gogitor.json`
5. Command-line flags

### Default config location

```text
~/.gogitor/config.json
```

Logs are written to:

```text
~/.gogitor/logs/gogitor_YYYY-MM-DD.log
```

### Example `config.json`

```json
{
  "provider": "ollama",
  "model": "gemma3:4b",
  "api_key": "",
  "ollama_url": "http://localhost:11434",
  "log_level": "info",
  "debug_mode": false,
  "dry_run": false,
  "llm_timeout": 300,
  "max_iterations": 5,
  "auto_git_commit": true,
  "git_auto_init": true,
  "multi_agent_enabled": true,
  "raw_output": false
}
```

### Provider formats

| Provider value | Meaning |
|---|---|
| `ollama` | Use local Ollama |
| `http://host:11434` | Use an Ollama-compatible server by URL |
| `https://host:11434` | Use an Ollama-compatible server by URL |
| `openai+https://api.example.com/v1` | Use an OpenAI-compatible API |
| `openai-compatible+http://localhost:8000/v1` | Use an OpenAI-compatible API |

### Environment variables

| Variable | Description |
|---|---|
| `GOGITOR_PROVIDER` | Default provider |
| `GOGITOR_MODEL` | Default model |
| `GOGITOR_API_KEY` | API key |
| `OPENAI_API_KEY` | Fallback API key for OpenAI-compatible providers |
| `GOGITOR_OLLAMA_URL` | Ollama URL |
| `OLLAMA_HOST` | Fallback Ollama host if URL is empty |
| `GOGITOR_LOG_LEVEL` | Log level |
| `GOGITOR_DEBUG` | Enable debug mode |
| `GOGITOR_DRY_RUN` | Enable dry-run by default |
| `GOGITOR_RAW` | Enable raw output by default |
| `GOGITOR_LLM_TIMEOUT` | LLM timeout in seconds |
| `GOGITOR_MAX_ITERATIONS` | Maximum code-fix iterations |
| `GOGITOR_AUTO_GIT_COMMIT` | Enable/disable automatic Git commit |
| `GOGITOR_GIT_AUTO_INIT` | Enable/disable automatic Git init |
| `GOGITOR_MULTI_AGENT` | Enable/disable multi-agent planning |

### Local project config

You can place `.gogitor.json` in the project root:

```json
{
  "provider": "ollama",
  "model": "gemma3:4b",
  "auto_git_commit": false,
  "dry_run": false
}
```

---

## Examples

### Ollama TUI

```bash
./gogitor tui --provider ollama --model gemma3:4b
```

### Remote Ollama-compatible server

```bash
./gogitor tui --provider http://192.168.1.10:11434 --model gemma3:4b
```

### OpenAI-compatible API

```bash
./gogitor ask "explain generics in Go" \
  --provider openai+https://api.example.com/v1 \
  --model gpt-4o-mini \
  --key sk-...
```

### Code generation

```bash
./gogitor code "create a REST API with health and version endpoints"
```

### Dry run

```bash
./gogitor code "refactor main.go" --dry-run
```

### Skip tests

```bash
./gogitor code "add logging" --no-tests
```

### Disable auto commit

```bash
./gogitor code "split code into files" --no-commit
```

### Task file

```bash
./gogitor task ./tasks/feature.txt
```

### Force code mode for task file

```bash
./gogitor file ./tasks/refactor.md --code
```

### JSON output

```bash
./gogitor test --json
```

### Raw output for pipes

```bash
echo "write hello world in Go" | ./gogitor code --raw > main.go
```

```bash
./gogitor ask "explain context.Context" --raw
```

### If task starts with a dash

Use `--` to separate flags from task text:

```bash
./gogitor code -- --fix something
```

---

## How code generation works

1. Gogitor receives your task.
2. It builds project context from existing files.
3. It asks the LLM to return either:
   - full file blocks:
     ```text
     --- File: path/to/file.go ---
     <file content>
     ```
   - or patch blocks:
     ```text
     --- Patch: path/to/file.go ---
     <<<<<<< SEARCH
     exact existing code
     =======
     new code
     >>>>>>> REPLACE
     ```
4. Gogitor copies your project into a temporary sandbox.
5. It applies the proposed changes there.
6. It runs:
   - `go build`
   - `go test`
7. If validation succeeds, changes are copied back to your project.
8. If enabled, Gogitor creates a Git commit.

If patch mode fails, Gogitor can fall back to full-file rewrite mode.

---

## Safety

Gogitor tries to be safe, but it still executes generated code and modifies files on your machine.

Recommendations:

- Use Git.
- Run with `--dry-run` when unsure.
- Review generated patches before committing.
- Use trusted models and endpoints.
- Do not run untrusted task files without review.

The sandbox used by Gogitor is a temporary directory, not a full container or VM. It helps prevent accidental broken builds in your working tree, but it is not a complete isolation mechanism.

---

## Project structure

```text
main.go                  CLI entry point
internal/app             application orchestration
internal/codegen         parsing of LLM file/patch output
internal/config          configuration loading and validation
internal/domain          shared domain types
internal/git             git operations
internal/llm             LLM client
internal/prompts         prompt builders
internal/runner          go build/test/run execution
internal/search          web search
internal/security        path safety helpers
internal/workspace       project files and sandbox operations
internal/ui/cli          CLI interface
internal/ui/tui          TUI interface
```

---

## Troubleshooting

### `unsupported provider`

Use one of the supported provider formats:

```bash
--provider ollama
--provider http://localhost:11434
--provider openai+https://api.example.com/v1
--provider openai-compatible+http://localhost:8000/v1
```

### Ollama is not reachable

Check that Ollama is running:

```bash
ollama serve
```

Or specify a custom URL:

```bash
./gogitor tui --provider http://127.0.0.1:11434 --model gemma3:4b
```

### Build fails

Make sure the project itself builds:

```bash
go build ./...
```

Then try Gogitor again.

### Tests fail after code generation

Run with more details:

```bash
./gogitor test --json
```

Or temporarily skip tests if you only want code output:

```bash
./gogitor code "task" --no-tests
```

---

## License

Add a `LICENSE` file to your repository if you plan to publish this project.