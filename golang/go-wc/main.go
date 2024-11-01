package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	seen := make(map[string]bool)
	defer fmt.Println("seen lines:\n", seen)

	done := make(chan bool)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("got: ", line)
			if !seen[line] {
				seen[line] = true
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("scanner error: ", err)
		}

		close(done)
	}()

	select {
	case sig := <-signals:
		fmt.Printf("received signal: %v\n", sig)
	case <-done:
		fmt.Println("finished normally")
	}
}
