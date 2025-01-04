package main

import (
	"fmt"
	"strings"
	"time"
)

func spinner() {
	for {
		for i, c := range "|/-\\" {
			fmt.Printf("\r%c Working on it%s%s", c, strings.Repeat(".", i), strings.Repeat(" ", 3-i))
			time.Sleep(time.Millisecond * 250)
		}
	}
}

func awfullyLongCalculation() {
	time.Sleep(60 * time.Second)
}

func main() {
	go spinner()
	awfullyLongCalculation()
}
