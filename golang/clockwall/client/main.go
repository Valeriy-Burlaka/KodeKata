package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	maxClocks        = 6
	clockWidth       = 16
	interClockMargin = 3
)

type Clock struct {
	City string
	Port int64
	Time chan string
}

func parse(args []string) (res []Clock, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("clockwall needs at least 1  max (got %d)", len(args))
	} else if len(args) > maxClocks {
		return nil, fmt.Errorf("clockwall can fit %d clocks max (got %d)", maxClocks, len(args))
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		parsed := strings.Split(args[i], "=")
		if len(parsed) != 2 {
			return nil, fmt.Errorf("invalid clock arg (want: CityName=port; got: %q)", arg)
		}
		maybeCity := parsed[0]
		maybePort := parsed[1]
		city := strings.TrimSpace(maybeCity)
		if city == "" {
			return nil, fmt.Errorf("city name %q in arg %q is invalid", maybeCity, arg)
		}
		port, e := strconv.Atoi(maybePort)
		if e != nil || port == 0 {
			return nil, fmt.Errorf("port number %q in arg %q is invalid", maybePort, arg)
		}

		res = append(res, Clock{City: city, Port: int64(port)})
	}

	return res, nil
}

func startClock(c *Clock) error {
	// TODO: should connect to a TCP time server
	return fmt.Errorf("not implemented")
}

func startClockMock(c *Clock) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	i := 0
	for range t.C {
		c.Time <- fmt.Sprintf("XX:XX:%02d", i%60)
		i++
	}
}

func clockTitle(c Clock) string {
	pad := float64(clockWidth-len(c.City)) / 2
	padLeft := int(math.Floor(pad))
	padRight := int(math.Ceil(pad))

	return fmt.Sprintf("%s%s%s",
		strings.Repeat(" ", padLeft),
		c.City,
		strings.Repeat(" ", padRight))
}

func buildTitleRow(clocks []Clock, margin string) []byte {
	titles := make([]string, len(clocks))
	for i := range clocks {
		titles[i] = clockTitle(clocks[i])
	}
	row := strings.Join(titles, margin)

	return []byte(row)
}

func buildDisplayRow(clocks []Clock, part string, margin string) []byte {
	parts := make([]string, len(clocks))
	for i := range clocks {
		parts[i] = part
	}
	row := strings.Join(parts, margin)

	return []byte(row)
}

func buildDisplay(clocks []Clock) ([][]byte, int) {
	endFrame := strings.Repeat("#", clockWidth)
	middleFrame := fmt.Sprintf("#%s#", strings.Repeat(" ", clockWidth-2))
	margin := strings.Repeat(" ", interClockMargin)

	display := make([][]byte, 6)

	display[0] = buildTitleRow(clocks, margin)
	display[1] = buildDisplayRow(clocks, endFrame, margin)
	display[2] = buildDisplayRow(clocks, middleFrame, margin)
	display[3] = buildDisplayRow(clocks, middleFrame, margin)
	display[4] = buildDisplayRow(clocks, middleFrame, margin)
	display[5] = buildDisplayRow(clocks, endFrame, margin)

	editableRow := 3

	return display, editableRow
}

func updateTime(c *Clock, row *[]byte, offset int) {
	for {
		t := <-c.Time
		padding := (clockWidth - len(t)) / 2
		for i := range t {
			(*row)[offset+padding+i] = t[i]
		}
	}
}

func main() {
	flag.Parse()

	clocks, err := parse(flag.Args())
	if err != nil {
		fmt.Printf("invalid args: %v\n", err)
		os.Exit(1)
	}

	clockDisplay, editableIndex := buildDisplay(clocks)

	for i := range clocks {
		c := &clocks[i]
		ch := make(chan string, 5)
		c.Time = ch

		go func() {
			err := startClock(c)
			if err != nil {
				// fmt.Printf("failed to start clock for %s: %v\n", c.City, err) - exit?
				startClockMock(c)
			}
		}()
		go updateTime(c, &clockDisplay[editableIndex], i*(clockWidth+interClockMargin))
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	// possible problem (DK how to solve yet) — synchronization between each clock timer and the main display timer.
	// - Is it possible that full re-draw happens when some clock is between time ticks?
	//   This would mean that some clocks may "lag" behind their peers.
	for range t.C {
		for _, r := range clockDisplay {
			fmt.Fprintln(os.Stdout, string(r))
		}
		// \033[nD - Move cursor back n characters
		// \033[nC - Move cursor forward n characters
		// \033[nA - Move cursor up n lines
		// \033[nB - Move cursor down n lines
		// \033[H - Move cursor to home position (0,0)
		// \033[n;mH - Move cursor to position (n,m)
		fmt.Print("\033[6A")
	}
}
