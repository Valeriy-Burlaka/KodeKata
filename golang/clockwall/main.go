package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	zones := []string{
		"Africa/Lagos",
		"America/Juneau",
		"America/Panama",
		"Asia/Riyadh",
		"Asia/Manila",
		"Australia/Melbourne",
		"Europe/Dublin",
		"Europe/Kyiv",
		"Europe/Warsaw",
	}

	for _, z := range zones {
		loc, err := time.LoadLocation(z)
		if err != nil {
			log.Fatalf("Load location failed: %v", err)
		}

		t := time.Now().In(loc)
		zname, offset := t.Zone()
		fmt.Printf("Time in %s is %s (%s, %d)\n", z, t.Format(time.Kitchen), zname, offset)
	}
}
