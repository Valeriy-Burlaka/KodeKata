package main

import (
	"fmt"
	"strings"
	"time"
)

func spinner(done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println()
			return
		default:
			for i, c := range "|/-\\" {
				fmt.Printf("\r%c Working on it%s%s", c, strings.Repeat(".", i), strings.Repeat(" ", 3-i))
				time.Sleep(time.Millisecond * 250)
			}
		}
	}
}

func awfullyLongCalculation() {
	time.Sleep(10 * time.Second)
}

func main() {
	done := make(chan bool)
	go spinner(done)
	awfullyLongCalculation()
	done <- true
	fmt.Println("Done!")
}
