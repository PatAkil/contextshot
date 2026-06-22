// Package pipeline orchestrates the capture -> OCR -> clipboard -> notify flow.
package pipeline

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/PatAkil/contextshot/internal/capture"
	"github.com/PatAkil/contextshot/internal/clip"
	"github.com/PatAkil/contextshot/internal/notify"
	"github.com/PatAkil/contextshot/internal/ocr"
)

// CaptureOnce runs one full capture: select a region, OCR it, copy the text to
// the clipboard, and show a confirmation. A user cancel is a clean no-op.
func CaptureOnce() error {
	path, err := capture.Region()
	if err != nil {
		if errors.Is(err, capture.ErrCancelled) {
			return nil
		}
		return err
	}
	defer func() {
		if rmErr := os.Remove(path); rmErr != nil && !os.IsNotExist(rmErr) {
			fmt.Fprintf(os.Stderr, "contextshot: removing temp capture: %v\n", rmErr)
		}
	}()

	text, err := ocr.Recognize()
	if err != nil {
		if errors.Is(err, ocr.ErrNoText) {
			return notify.Show("No text detected")
		}
		return err
	}

	if err := clip.Write(text); err != nil {
		return err
	}
	return notify.Show(notify.FormatMessage(utf8.RuneCountInString(text)))
}
