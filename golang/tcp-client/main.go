package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
		msg := string(buf[:n])

		if n > 0 {
			log.Printf("Received %q (%d bytes)", msg, n)
			if msg == "%QUIT%" {
				log.Println("server ended communication with %QUIT% msg")
				return nil
			}
		}

		if err == io.EOF {
			log.Println("Server ended communication")
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
