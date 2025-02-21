package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var port uint
var tz string
var loc *time.Location

const timecastInterval = 1 * time.Second
const timecastWriteTimeout = time.Duration(0.75 * float64(timecastInterval))

func init() {
	flag.UintVar(&port, "port", 8010, "Clock server port")
	flag.StringVar(&tz, "tz", "Europe/Kyiv", "Server time zone, as a valid IANA Time Zone name")

	flag.Parse()

	var err error
	loc, err = time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Failed to load location from tz name %q: %v", tz, err)
	}
}

func main() {
	fmt.Printf("Streaming local time in %q on port: %d\n", tz, port)

	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		log.Printf("accepted connection from %s", conn.LocalAddr().String())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	if err := startBroadcasting(conn); err != nil {
		log.Printf("Failed to handle connection: %v", err)
	}
}

func startBroadcasting(conn net.Conn) error {
	ticker := time.NewTicker(timecastInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := writeMsg(conn); err != nil {
			return fmt.Errorf("handleConn: %w", err)
		}
	}

	return nil
}

func writeMsg(conn net.Conn) error {
	t := time.Now().In(loc)

	if err := conn.SetWriteDeadline(time.Now().Add(timecastWriteTimeout)); err != nil {
		return fmt.Errorf("failed to set deadline: %w", err)
	}

	msg := fmt.Sprintf("%s: %s\n", tz, t.Format(time.TimeOnly))
	if _, err := io.WriteString(conn, msg); err != nil {
		return fmt.Errorf("failed to write to conn: %w", err)
	}

	return nil
}
