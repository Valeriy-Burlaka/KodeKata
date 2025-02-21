package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

var port uint
var tz string
var loc *time.Location

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

func handleConn(conn *net.Conn) (n int64, err error) {
	defer (*conn).Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		t := time.Now().In(loc)
		msg := fmt.Sprintf("%s: %s\n", tz, t.Format(time.TimeOnly))
		r := strings.NewReader(msg)

		err = (*conn).SetWriteDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Printf("error setting conn write timeout: %v", err)
			break
		}

		i, err := io.Copy(*conn, r)
		n += i
		if err != nil {
			log.Printf("error writing to conn (%d bytes sent): %v", i, err)
			break
		}
	}

	return n, err
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
			log.Printf("error accepting connection: %v", err)
			return
		}
		log.Printf("accepted connection from %s", conn.LocalAddr().String())

		go func(c net.Conn) {
			nbytes, err := handleConn(&conn)
			if err != nil {
				log.Printf("error handling connection: %v", err)
			}
			log.Printf("wrote %d bytes to connection", nbytes)
		}(conn)
	}
}
