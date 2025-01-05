package main

import (
	"fmt"
	"strings"
	"time"
)

func spinner(done chan struct{}) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	progress := 0
	spinnerChars := []rune{'|', '/', '-', '\\'}
	for {
		select {
		case <-done:
			fmt.Println()
			return
		case <-ticker.C:
			charIndex := progress % len(spinnerChars)
			currentChar := spinnerChars[charIndex]
			fmt.Printf("\r%c Working on it%s%s",
				currentChar,
				strings.Repeat(".", charIndex),
				strings.Repeat(" ", 3-charIndex))
			progress++
		}
	}
}

func awfullyLongCalculation(duration time.Duration) {
	time.Sleep(duration)
}

func main() {
	done := make(chan struct{})
	go spinner(done)
	awfullyLongCalculation(5250 * time.Millisecond)
	done <- struct{}{}
	fmt.Println("Done!")
}
