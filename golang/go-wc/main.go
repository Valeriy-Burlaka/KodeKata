package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	seen := make(map[string]bool)
	defer func() {
		if len(seen) == 0 {
			fmt.Println("no lines processed")
		} else {
			fmt.Println("processed lines:")
			for line := range seen {
				fmt.Printf(" - %s\n", line)
			}
		}
	}()

	done := make(chan bool)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			// fmt.Println("got: ", line)
			if !seen[line] {
				seen[line] = true
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("scanner error: %v", err)
		}

		close(done)
	}()

	select {
	case sig := <-signals:
		fmt.Printf("interrupted by signal: %v\n", sig)
	case <-done:
		fmt.Println("program finished normally")
	}
}
