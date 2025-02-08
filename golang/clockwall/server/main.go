package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var port int64
var tz string
var loc *time.Location

func init() {
	flag.Int64Var(&port, "port", 8010, "clock server port")

	flag.Parse()

	tz = os.Getenv("TZ")
	if tz == "" {
		log.Fatal("TZ environment variable is not set")
	}

	var err error
	loc, err = time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Failed to load location from tz name %q: %v", tz, err)
	}
}

func handleConn(conn *net.Conn) (n int64, err error) {
	defer (*conn).Close()

	t := time.Now().In(loc)
	msg := fmt.Sprintf("%s: %s", tz, t.Format(time.TimeOnly))
	r := strings.NewReader(msg)
	i, _ := io.Copy(*conn, r)

	n += i

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
			continue
		}
		handleConn(&conn)
	}
}
