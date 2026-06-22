// Command contextshot is a menu-bar daemon: press Cmd+Shift+2 (or use the menu)
// to select a screen region, OCR it, and copy the text to the clipboard.
package main

import (
	"log"

	"github.com/getlantern/systray"

	"github.com/PatAkil/contextshot/internal/hotkey"
	"github.com/PatAkil/contextshot/internal/notify"
	"github.com/PatAkil/contextshot/internal/pipeline"
)

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetTitle("📸")
	systray.SetTooltip("Contextshot — capture screen text")

	mCapture := systray.AddMenuItem("Capture now", "Capture text from a screen region")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Contextshot")

	if err := hotkey.Listen(runCapture); err != nil {
		log.Printf("contextshot: global hotkey unavailable: %v", err)
	}

	go func() {
		for {
			select {
			case <-mCapture.ClickedCh:
				runCapture()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

// runCapture runs one capture on a background goroutine so the menu and hotkey
// listeners are never blocked by the interactive selector.
func runCapture() {
	go func() {
		if err := pipeline.CaptureOnce(); err != nil {
			log.Printf("contextshot: capture failed: %v", err)
			if nErr := notify.Show("Capture failed"); nErr != nil {
				log.Printf("contextshot: showing failure notification: %v", nErr)
			}
		}
	}()
}
