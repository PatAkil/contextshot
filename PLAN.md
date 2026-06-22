# Contextshot — Option A Implementation Plan (Go)

**Goal:** A Go menu-bar daemon. Press a global hotkey → native region select → OCR via macOS Shortcuts → text lands on clipboard → brief confirmation. No hand-written cgo; no Xcode project.

**Name:** Contextshot — a "screenshot" that captures *text* (not pixels), built to paste screen content as context into agents.

---

## Stack & dependencies

| Concern | Tool | Why |
|---|---|---|
| Region capture | `screencapture` (built into macOS) | Native crosshair UI for free; matches the native UX |
| OCR | macOS **Shortcuts** CLI (`shortcuts run`) | "Extract Text from Image" action = Apple Vision, local, no API |
| Clipboard | `golang.design/x/clipboard` (or `pbcopy`) | Direct write; clean package |
| Global hotkey | `golang.design/x/hotkey` | Maintained; hides the cgo |
| Menu bar | `github.com/getlantern/systray` | Menu-bar icon + menu in pure-ish Go |
| Notification | `osascript -e 'display notification ...'` | Zero-dep confirmation toast |

One-time manual setup: create a Shortcut (see Phase 1) — the only non-code prerequisite.

---

## Project layout
```
contextshot/
├── go.mod
├── main.go              # systray lifecycle, wiring
├── internal/
│   ├── capture/         # screencapture wrapper → temp PNG path
│   ├── ocr/             # shortcuts run wrapper → string
│   ├── clip/            # clipboard write
│   ├── hotkey/          # hotkey registration → fires capture flow
│   └── notify/          # osascript toast
└── assets/icon.png      # menu-bar template icon
```

---

## Phase 0 — Prove the risky core (do this first, ~30 min)
Before any UI or hotkeys, validate the capture→OCR→clipboard chain from a plain `main.go`:
1. `screencapture -i -x /tmp/contextshot.png`
2. `shortcuts run "Contextshot OCR" -i /tmp/contextshot.png` → capture stdout
3. write stdout to clipboard
4. paste somewhere, confirm the text is right

If this works end-to-end, everything else is just glue. **This is the make-or-break step** — the Shortcuts CLI output format is the one unknown.

## Phase 1 — The OCR Shortcut (manual, one-time)
In the **Shortcuts** app, create a shortcut named `Contextshot OCR`:
- Input: **Shortcut Input** (set to accept Images/Files)
- Action: **Extract Text from Image**
- Action: output the text (return on stdout) — preferred over copy-to-clipboard so Go can add a preamble later

Test from terminal:
```
shortcuts run "Contextshot OCR" -i /tmp/contextshot.png
```

## Phase 2 — `capture` package
```go
func Region() (string, error) {
    path := "/tmp/contextshot.png"
    cmd := exec.Command("screencapture", "-i", "-x", path)
    if err := cmd.Run(); err != nil { return "", err }
    // -i lets user hit Esc → file won't exist; treat as cancel
    if _, err := os.Stat(path); err != nil {
        return "", ErrCancelled
    }
    return path, nil
}
```
Handle the **cancel case** (user presses Esc) cleanly — no error toast, just abort.

## Phase 3 — `ocr` package
```go
func Recognize(imagePath string) (string, error) {
    out, err := exec.Command("shortcuts", "run",
        "Contextshot OCR", "-i", imagePath).Output()
    if err != nil { return "", err }
    return strings.TrimSpace(string(out)), nil
}
```
Edge: empty result (no text found) → toast "No text detected", don't clobber clipboard.

## Phase 4 — `clip` + `notify`
```go
clipboard.Write(clipboard.FmtText, []byte(text))   // golang.design/x/clipboard

// notify
exec.Command("osascript", "-e",
    `display notification "Copied 142 chars" with title "Contextshot"`).Run()
```

## Phase 5 — `hotkey` + the flow
Register e.g. `Cmd+Shift+2`. On fire, run the pipeline on a goroutine (don't block the hotkey listener):
```
hotkey fired
  → capture.Region()      (ErrCancelled → silently stop)
  → ocr.Recognize(path)   (empty → "no text" toast)
  → clip.Write(text)
  → notify("Copied N chars")
  → os.Remove(tmp)
```

## Phase 6 — `systray` shell
- systray runs without a dock icon by default.
- Menu items: **Capture now**, **Preferences** (later), **Quit**.
- Register the hotkey inside `systray.onReady`.

## Phase 7 — Packaging
- `go build` produces a binary. For a real `.app`, wrap it: minimal `Contextshot.app/Contents/MacOS/contextshot` + `Info.plist` with `LSUIElement=true`.
- First run triggers the **Screen Recording** permission prompt (required for `screencapture` content on modern macOS) — document this.

---

## Risks / watch-items
1. **Shortcuts CLI output** — `shortcuts run` stdout behavior can be finicky across macOS versions. Phase 0 de-risks this; fallback is the Shortcut copies to clipboard itself and Go reads the clipboard back.
2. **Permissions** — Screen Recording prompt is unavoidable and only appears on first capture.
3. **`golang.design/x/hotkey`** uses cgo under the hood, so building needs Xcode command-line tools (`xcode-select --install`). You won't write C, but it's compiled.
4. **Latency** — spawning `shortcuts` per capture adds ~0.5–1s. Fine for personal use; if it annoys you, move OCR into a Swift helper binary later.

---

## Build order
**Phase 0 (validate core) → 2 → 3 → 4 → 5 → 6 → 7.** Get a working terminal script before touching hotkeys or the menu bar.

### Stretch (v2)
- Preamble toggle: prepend `Context from my screen:\n\n` before clipboard write — the original "paste into agents" use case.
- Replace the Shortcuts shell-out with a tiny embedded Swift OCR binary for speed + no manual Shortcut setup.
- Capture history (last N grabs).
