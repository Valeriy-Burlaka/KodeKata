package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var host string
var port int

const (
	bufSize     = 4096
	connTimeout = 5 * time.Second
	maxAttempts = 3
	readTimeout = 2 * time.Second
)

const EOM = "%QUIT%"

func init() {
	flag.StringVar(&host, "h", "localhost", "Specify host")
	flag.IntVar(&port, "p", 8000, "Specify port")

	flag.Parse()
}

func read(conn net.Conn) error {
	attempts := maxAttempts
	buf := make([]byte, bufSize)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
			return fmt.Errorf("set deadline error: %w", err)
		}

		n, err := conn.Read(buf)
		msg := strings.TrimSpace(string(buf[:n]))
		if n > 0 {
			var end bool
			if strings.Contains(msg, EOM) {
				end = true
				msg = strings.TrimSpace(strings.Replace(msg, EOM, "", -1))
			}
			log.Printf("Received %q (%d bytes)", msg, n)

			if end {
				log.Printf("The server ended communication with %q msg", EOM)
				return nil
			}
		}

		if err == io.EOF {
			log.Println("The server ended communication")
			return nil
		}
		if errors.Is(err, os.ErrDeadlineExceeded) {
			if attempts <= 0 {
				return fmt.Errorf("max read attempts exceeded")
			}
			log.Printf("Read timeout exceeded (%d attempts left)", attempts)
			attempts--
			continue
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
	}
}

func main() {
	fmt.Printf("host = %s\n", host)
	fmt.Printf("port = %d\n", port)
	connStr := fmt.Sprintf("%v:%v", host, port)
	fmt.Println("connecting to addr ", connStr)

	conn, err := net.DialTimeout("tcp", connStr, connTimeout)
	if err != nil {
		log.Fatalf("Error connecting to %q: %v", connStr, err)
	}
	defer conn.Close()

	if err := read(conn); err != nil {
		log.Fatalf("Error reading from connection: %v", err)
	}
}
