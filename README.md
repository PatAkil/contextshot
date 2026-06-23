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

1. Click **+** to create a new shortcut and rename it (top center) to
   `Contextshot OCR`.
2. Add the action **Extract Text from Image**.
3. Make the shortcut accept an image as input: open the shortcut's **ⓘ**
   (Details) panel and turn on **Show in Share Sheet**, then drag the
   **Shortcut Input** variable into the action's "Image" field (so the image
   passed with `-i` flows into the OCR action).
4. The shortcut's final output must be the extracted text — Contextshot reads it
   from the shortcut's file output (`-o`). No "Copy to Clipboard" action needed.

Verify it from the terminal (the output should be the text in the image):

```bash
shortcuts run "Contextshot OCR" -i /path/to/some-image.png -o /tmp/out.txt && cat /tmp/out.txt
```

### 2. Grant Screen Recording

The **only** permission Contextshot needs is **Screen Recording** (so
`screencapture` can read the screen). The global hotkey uses Carbon and needs
**no** Accessibility or Input Monitoring permission.

On your first capture, macOS prompts "Contextshot would like to record this
screen" → **Open System Settings** → enable **Contextshot** under
**Privacy & Security → Screen Recording**, then let it **Quit & Reopen**.

> **If you run via `just run` instead of the app**, the permission attaches to
> your *terminal* app, not Contextshot — granting your terminal Screen Recording
> lets *anything* it launches read the screen. Prefer the packaged app below so
> the permission belongs to Contextshot itself.

## Build & run

Install dev tools and run the standalone menu-bar app (no Dock icon, via
`LSUIElement`) — this is the recommended way to run it:

```bash
just init
just package
open dist/Contextshot.app
```

`just package` ad-hoc code-signs the bundle with a stable identifier
(`com.patakil.contextshot`) so its code signature and app identity agree, which
is required for the Screen Recording grant to apply.

For quick dev iteration you can also run from source (`just run`), but a Dock
icon appears and Screen Recording attaches to your terminal (see note above).

### Re-granting Screen Recording after a rebuild

The bundle is **ad-hoc signed**, so every `just package` produces a new code
hash and macOS treats it as a different app — the Screen Recording toggle may
still show "on" but the app gets re-prompted (or capture silently fails). When
that happens, clear the stale grant and re-allow:

```bash
tccutil reset ScreenCapture com.patakil.contextshot
```

Then remove any stale **Contextshot** row from **Screen Recording** (select it,
click **−**), relaunch the app, trigger a capture, and accept the fresh prompt.

To avoid this entirely, sign the app with a **stable self-signed identity**
(Keychain Access → Certificate Assistant → *Create a Certificate* → type
**Code Signing**), then replace the `--sign -` in `scripts/package.sh` with your
certificate name. With a stable identity the Screen Recording grant survives
rebuilds.

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

- Press **Cmd+Shift+2** anywhere, or click the menu-bar **viewfinder icon →
  Capture now** (the icon is at the top-right of the screen; it's a menu-bar-only
  app with no Dock icon or window).
- Drag to select a region (press **Esc** to cancel — no error).
- The recognized text is on your clipboard; a toast confirms the character count.
- Quit from the menu-bar **viewfinder icon → Quit**.
- To launch at login: System Settings → General → Login Items → **+** →
  select `dist/Contextshot.app`.

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

## Acknowledgments

The project skeleton — the `justfile`, `just ci` validation pipeline, semgrep
rules, `arch-go.yml`, and `AGENTS.md` — was generated from the `go-cli-base`
blueprint of [**ai-guardrails**](https://github.com/florianbuetow/ai-guardrails)
by Florian Buetow. See that repository for its license and terms.

## License

[MIT](LICENSE) © 2026 Patrick Akil
