package spinner

import (
	"fmt"
	"time"
)

func Wait(done chan struct{}) {
	spinnerStages := []string{"▉", "▊", "▋", "▌", "▍", "▎", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}
	var passed time.Duration = 0
	for {
		for _, s := range spinnerStages {
			select {
			case <-done:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r%s %.2fs", s, passed.Seconds())
				passed += 100 * time.Millisecond
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
