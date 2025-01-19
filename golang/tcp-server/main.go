package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

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

		func() {
			ticker := time.NewTicker(1 * time.Second)

			defer conn.Close()
			defer ticker.Stop()

			for range ticker.C {
				err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				if err != nil {
					log.Printf("Error setting write deadline: %v", err)
					return
				}

				_, err = conn.Write([]byte(time.Now().Format(time.RFC1123)))
				if err != nil {
					log.Printf("Error writing to connection: %v", err)
					return
				}

			}
		}()

		log.Printf("Client %v disconnected", conn.RemoteAddr())
	}
}
