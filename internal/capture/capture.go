// Package capture wraps the macOS screencapture tool to grab a user-selected
// screen region as a PNG file.
package capture

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ImagePath is the fixed location of the captured region. It is a constant so
// the path passed to screencapture (and later to OCR) is never derived from
// external input.
const ImagePath = "/tmp/contextshot.png"

// ErrCancelled is returned when the user dismisses the region selector (Esc)
// without capturing anything. Callers should treat this as a clean abort.
var ErrCancelled = errors.New("capture cancelled")

// Region shows the native interactive region selector and returns the path to
// the captured PNG. If the user cancels, it returns ErrCancelled.
func Region() (string, error) {
	// Remove any leftover file so a stale capture can't masquerade as success
	// when the user cancels this run.
	if err := os.Remove(ImagePath); err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("clearing previous capture: %w", err)
	}

	// -i: interactive region selection, -x: no capture sound.
	cmd := exec.Command("screencapture", "-i", "-x", ImagePath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running screencapture: %w", err)
	}

	// On cancel, screencapture exits 0 but writes no file.
	if _, err := os.Stat(ImagePath); err != nil {
		if os.IsNotExist(err) {
			return "", ErrCancelled
		}
		return "", fmt.Errorf("checking capture output: %w", err)
	}
	return ImagePath, nil
}
