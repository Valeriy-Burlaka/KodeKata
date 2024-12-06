package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Event represents a server-sent event
type Event struct {
	Event string
	Data  string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/events", handleSSE)

	log.Printf("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	log.Printf("Client %v connected", r.RemoteAddr)
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()

	messageChan := make(chan Event)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Start a goroutine to send periodic events
	go func() {
		count := 0
		for {
			select {
			case <-ticker.C:
				count++
				messageChan <- Event{
					Event: "ping",
					Data:  fmt.Sprintf("Server time: %v (msg #%d)", time.Now().Format(time.RFC3339), count),
				}
			}
		}
	}()

	// Stream events to client
	for {
		select {
		case <-ctx.Done():
			log.Printf("Client %v disconnected (err=%v)", r.RemoteAddr, ctx.Err())
			return
		case msg := <-messageChan:
			// Format the event according to SSE specification
			fmt.Fprintf(w, "event: %s\n", msg.Event)
			fmt.Fprintf(w, "data: %s\n\n", msg.Data)
			w.(http.Flusher).Flush()
		}
	}
}
