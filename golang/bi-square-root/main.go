package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const eps = 0.001
const maxAttempts = 100

func main() {
	flag.Parse()

	maybeNum := flag.Arg(0)
	num, err := strconv.ParseFloat(maybeNum, 64)
	if err != nil {
		fmt.Printf("%q is not a valid number\n", maybeNum)
		os.Exit(1)
	}
	if num <= 0 {
		fmt.Printf("Must be a positive number, got %v\n", num)
		os.Exit(1)
	}

	fmt.Printf("Finding the square root of %v (error tolerance = %v)\n", num, eps)

	var attempt = 1
	var low float64 = 0
	var high float64 = num
	for {
		maybeAnsw := (low + high) / 2
		squared := maybeAnsw * maybeAnsw
		fmt.Printf("Attempt %d: Trying %v as the answer (%v squared)\n", attempt, maybeAnsw, squared)
		if squared >= num-eps && squared <= num+eps {
			fmt.Printf("Found possible answer: %v\n", maybeAnsw)
			break
		}
		if squared > num {
			high = maybeAnsw
		} else {
			low = maybeAnsw
		}

		attempt++
		if attempt > maxAttempts {
			fmt.Printf("Error: Could not find the answer after %d attempts\n", attempt)
			break
		}
	}
}
