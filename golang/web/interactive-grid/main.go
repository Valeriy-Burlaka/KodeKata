package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
)

// Event represents a server-sent event
type Event struct {
	Event string
	Data  string
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

func handleSpaceNew(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := generateID()
		if err != nil {
			log.Printf("failed to generate ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		space := &Space{
			ID:        id,
			Rows:      10,
			Cols:      10,
			Enabled:   make([]string, 0),
			CreatedAt: time.Now(),
		}

		if err := store.SaveSpace(space); err != nil {
			log.Printf("failed to save space: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, path.Join("/space", id), http.StatusSeeOther)
	}
}

var indexTmpl *template.Template
var spaceTmpl *template.Template

func init() {
	indexTmpl = template.Must(template.ParseFiles("index.html"))
	spaceTmpl = template.Must(template.ParseFiles("space.html"))
}

type IndexPageData struct {
	Spaces []*Space
}

func handleIndex(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.mu.RLock()
		spaces := make([]*Space, 0, len(store.spaces))
		for _, space := range store.spaces {
			spaces = append(spaces, space)
		}
		store.mu.RUnlock()

		data := IndexPageData{Spaces: spaces}
		if err := indexTmpl.Execute(w, data); err != nil {
			log.Printf("failed to execute template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

type SpacePageData struct {
	Space *Space
}

func handleSpace(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		spaceID := vars["id"]

		// Get space data from store
		space, err := store.GetSpace(spaceID)
		if err != nil {
			log.Printf("Error getting space %s: %v", spaceID, err)
			http.Error(w, "Space not found", http.StatusNotFound)
			return
		}

		// Prepare template data
		data := SpacePageData{
			Space: space,
		}

		// Set content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render template
		if err := spaceTmpl.Execute(w, data); err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	store, err := NewStore("spaces.json")
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", handleIndex(store)).Methods("GET")

	r.HandleFunc("/spaces/new", handleSpaceNew(store)).Methods("POST")
	r.HandleFunc("/space/{id}", handleSpace(store)).Methods("GET")

	http.HandleFunc("/events", handleSSE)

	addr := fmt.Sprintf("%s:%d", "localhost", 8080)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
