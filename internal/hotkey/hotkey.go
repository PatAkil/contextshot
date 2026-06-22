// Package hotkey registers the global capture hotkey (Cmd+Shift+2) and invokes a
// callback each time it is pressed.
package hotkey

import (
	"fmt"

	"golang.design/x/hotkey"
)

// Handler is called once per hotkey press.
type Handler func()

// Listen registers the global hotkey and starts a goroutine that calls onPress
// for every keypress. The goroutine holds the only reference to the registered
// hotkey, keeping it alive for the lifetime of the program.
func Listen(onPress Handler) error {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.Key2)
	if err := hk.Register(); err != nil {
		return fmt.Errorf("registering Cmd+Shift+2 hotkey: %w", err)
	}
	go func() {
		for range hk.Keydown() {
			onPress()
		}
	}()
	return nil
}
