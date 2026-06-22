// Package clip writes text to the macOS clipboard via pbcopy.
package clip

import (
	"fmt"
	"os/exec"
	"strings"
)

// Write copies text to the system clipboard.
func Write(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("writing to clipboard: %w", err)
	}
	return nil
}
