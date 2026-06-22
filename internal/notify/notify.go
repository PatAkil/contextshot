// Package notify shows a brief macOS notification via osascript.
package notify

import (
	"fmt"
	"os"
	"os/exec"
)

// messageScript reads the notification text from an environment variable so the
// osascript argument is a constant. AppleScript's `system attribute` returns the
// value of the named environment variable.
const messageScript = `display notification (system attribute "CONTEXTSHOT_MSG") with title "Contextshot"`

// FormatMessage builds the confirmation text for a successful copy of n
// characters.
func FormatMessage(n int) string {
	if n == 1 {
		return "Copied 1 character"
	}
	return fmt.Sprintf("Copied %d characters", n)
}

// Show displays a notification toast with the given message.
func Show(message string) error {
	cmd := exec.Command("osascript", "-e", messageScript)
	cmd.Env = append(os.Environ(), "CONTEXTSHOT_MSG="+message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("showing notification: %w", err)
	}
	return nil
}
