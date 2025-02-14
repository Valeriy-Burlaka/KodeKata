package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const EOM = "%QUIT%"

func write(conn net.Conn) error {
	maxMsgs := 10
	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()
	defer conn.Close()
	defer log.Printf("Client %v disconnected", conn.RemoteAddr())

	i := 1
	for range ticker.C {
		err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			return fmt.Errorf("error setting write deadline: %w", err)
		}

		_, err = conn.Write([]byte(time.Now().Format(time.RFC1123) + "\n"))
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		i++
		if i > maxMsgs {
			_, err := conn.Write([]byte(EOM + "\n"))
			if err != nil {
				return fmt.Errorf("write error: %w", err)
			}
			return nil
		}
	}

	return nil
}

func main() {
	fmt.Println("Starting TCP server on port 8000")
	server, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("New client connected from %v", conn.RemoteAddr())

		go func() {
			if err := write(conn); err != nil {
				log.Printf("Error writing to connection: %v", err)
			}
		}()
	}
}
