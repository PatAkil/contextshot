# contextshot

A macOS menu-bar daemon that captures *text* from your screen. Press a global
hotkey, drag a region, and the text inside it is OCR'd and placed on your
clipboard — ready to paste as context into an agent. A "screenshot" for text.

## How it works

```
Cmd+Shift+2 (or menu "Capture now")
  → screencapture -i      native region selector (Esc cancels)
  → "Contextshot OCR"     Apple Vision OCR via the Shortcuts CLI (local)
  → pbcopy                text lands on the clipboard
  → osascript toast       "Copied N characters"
```

All processing is local — no network, no API keys. OCR is done by Apple Vision
through a one-time Shortcut you create yourself.

## Prerequisites

- **macOS 13+** with **Xcode command-line tools** (`xcode-select --install`) —
  the hotkey and menu-bar libraries use cgo.
- **Go 1.24+** (the module pins `toolchain go1.26.4`, auto-downloaded).
- For development/CI only: **just**, **codespell**, **semgrep**
  (`brew install just codespell semgrep`).

## One-time setup

### 1. Create the OCR Shortcut

In the **Shortcuts** app, create a shortcut named exactly `Contextshot OCR`:

1. Add the action **Extract Text from Image**.
2. Set its input to the **Shortcut Input** (accept Images/Files).
3. The shortcut's final output must be the extracted text (Contextshot reads it
   from the shortcut's file output — no "copy to clipboard" action needed).

Verify it from the terminal:

```bash
shortcuts run "Contextshot OCR" -i /path/to/some-image.png -o /tmp/out.txt && cat /tmp/out.txt
```

### 2. Grant permissions (prompted on first use)

- **Screen Recording** — required for `screencapture` to read screen content.
- **Accessibility** — required for the global hotkey (Cmd+Shift+2).

System Settings → Privacy & Security → enable Contextshot (or your terminal, if
running via `just run`) under both sections, then restart the app.

## Build & run

Install dev tools and run from source (a Dock icon appears in this mode):

```bash
just init
just run
```

Build the standalone menu-bar app (no Dock icon, via `LSUIElement`):

```bash
just package
open dist/Contextshot.app
```

## Repository Structure

```
contextshot/
├── go.mod / go.sum         # Go module + checksums
├── justfile                # Task runner (init, run, package, ci, ...)
├── scripts/package.sh      # Builds Contextshot.app
├── cmd/
│   └── main.go             # systray menu-bar lifecycle + hotkey wiring
├── internal/
│   ├── capture/            # screencapture -i wrapper (Esc → ErrCancelled)
│   ├── ocr/                # "Contextshot OCR" Shortcut wrapper
│   ├── clip/               # pbcopy clipboard write
│   ├── notify/             # osascript notification toast
│   ├── hotkey/             # global Cmd+Shift+2 registration
│   └── pipeline/           # capture → OCR → clipboard → notify flow
├── AGENTS.md               # AI agent development rules (CLAUDE.md mirrors this)
├── config/                 # semgrep + codespell configs
└── reports/                # generated coverage/security reports (not in git)
```

## Usage

- Press **Cmd+Shift+2** anywhere, or click the menu-bar **📸 → Capture now**.
- Drag to select a region (press **Esc** to cancel — no error).
- The recognized text is on your clipboard; a toast confirms the character count.
- Quit from the menu-bar **📸 → Quit**.

See all developer commands with `just help`.

## Development

### Available Commands

- `just init` - Initialize development environment
- `just run` - Run the main application
- `just destroy` - Remove build artifacts and caches
- `just help` - Show available commands

### Code Quality

- `just code-style` - Check code formatting (read-only)
- `just code-format` - Auto-fix code formatting
- `just code-typecheck` - Run go vet
- `just code-lspchecks` - Run staticcheck
- `just code-security` - Run security checks (gosec)
- `just code-deptry` - Check dependency hygiene
- `just code-spell` - Check spelling
- `just code-audit` - Scan for vulnerabilities (govulncheck)
- `just code-semgrep` - Run custom static analysis
- `just lint` - Run golangci-lint (meta-linter)
- `just code-architecture` - Run architecture checks (arch-go)

### Testing

- `just test` - Run unit tests
- `just test-coverage` - Run tests with coverage (0% threshold)

### CI

- `just ci` - Run all validation checks (verbose)
- `just ci-quiet` - Run all checks (silent, fail-fast)

The CI pipeline runs the following steps in order:
1. `init` - Initialize environment
2. `code-format` - Auto-format code
3. `code-style` - Verify formatting
4. `code-typecheck` - Type checking (go vet)
5. `code-security` - Security scan (gosec)
6. `code-deptry` - Dependency hygiene
7. `code-spell` - Spell checking
8. `code-semgrep` - Custom static analysis
9. `code-audit` - Vulnerability scanning (govulncheck)
10. `test` - Unit tests
11. `code-architecture` - Architecture checks (arch-go)
12. `lint` - golangci-lint
13. `code-lspchecks` - Strict static analysis (staticcheck)

## Project Rules

See [AGENTS.md](AGENTS.md) for detailed development guidelines including:
- Go execution rules
- Git commit guidelines (no AI attribution)
- Testing requirements
- Project structure conventions

## License

<!-- Add your license here -->
