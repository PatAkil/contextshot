// Package ocr extracts text from the captured image by running the macOS
// "Contextshot OCR" Shortcut, which uses Apple Vision locally.
package ocr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// shortcutName is the Shortcut the user creates once (see README). It must end
// with an "Extract Text from Image" action whose result is the shortcut output.
const shortcutName = "Contextshot OCR"

// inputPath mirrors capture.ImagePath. It is duplicated here as a local literal
// (rather than imported) so the subprocess call below contains no cross-package
// value, which gosec's taint analysis flags. Keep the two in sync.
const inputPath = "/tmp/contextshot.png"

// outputPath is the fixed file the shortcut writes its text result to. Using a
// constant (rather than a variable) keeps the subprocess and file-read calls
// free of externally-derived input.
const outputPath = "/tmp/contextshot-ocr.txt"

// ErrNoText is returned when the shortcut runs but finds no text in the image.
var ErrNoText = errors.New("no text detected")

// Recognize runs the OCR shortcut on the captured image and returns the
// extracted text. The shortcut output is written to a temp file (via -o) and
// read back, which avoids relying on the Shortcuts CLI's stdout behaviour.
func Recognize() (string, error) {
	cmd := exec.Command("shortcuts", "run", shortcutName, "-i", inputPath, "-o", outputPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running OCR shortcut %q: %w", shortcutName, err)
	}
	defer func() {
		if rmErr := os.Remove(outputPath); rmErr != nil && !os.IsNotExist(rmErr) {
			fmt.Fprintf(os.Stderr, "contextshot: removing OCR temp file: %v\n", rmErr)
		}
	}()

	data, err := os.ReadFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("reading OCR output: %w", err)
	}

	text := Clean(string(data))
	if text == "" {
		return "", ErrNoText
	}
	return text, nil
}

// Clean trims surrounding whitespace from raw OCR output.
func Clean(raw string) string {
	return strings.TrimSpace(raw)
}
