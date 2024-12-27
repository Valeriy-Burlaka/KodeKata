package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Space struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt"`
	ActiveCells int        `json:"activeCells"`
	Clients     int        `json:"clients"`
}

var spaceClients = make(map[string][]chan []byte) // Store SSE clients per space

func main() {
	// Initialize data files
	if _, err := os.Stat("spaces.json"); os.IsNotExist(err) {
		if err := os.WriteFile("spaces.json", []byte("[]"), 0644); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("passwords.json"); os.IsNotExist(err) {
		if err := os.WriteFile("passwords.json", []byte("{}"), 0644); err != nil {
			log.Fatal(err)
		}
	}

	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/spaces", handleSpaces)
	http.HandleFunc("/spaces/events", handleSpacesEvents)
	http.HandleFunc("/spaces/new", handleNewSpace)
	http.HandleFunc("/spaces/", handleSpace)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSpaces(w http.ResponseWriter, r *http.Request) {
	spaces, err := loadSpaces()
	if err != nil {
		http.Error(w, "Error loading spaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spaces); err != nil {
		http.Error(w, "Error encoding spaces", http.StatusInternalServerError)
		return
	}
}

func handleSpacesEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // For testing purposes

	// Create a channel for this client
	clientChan := make(chan []byte)

	// Add client channel to the global map
	spaceClients["spaces"] = append(spaceClients["spaces"], clientChan)

	// Remove client channel when the client disconnects
	defer func() {
		spaceClients["spaces"] = removeChannel(spaceClients["spaces"], clientChan)
		close(clientChan)
	}()

	for {
		select {
		case <-r.Context().Done():
			return // Client disconnected
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func removeChannel(channels []chan []byte, target chan []byte) []chan []byte {
	var newChannels []chan []byte
	for _, c := range channels {
		if c != target {
			newChannels = append(newChannels, c)
		}
	}
	return newChannels
}

func broadcastSpaceUpdate(spaceID string, spaceData []byte) {
	if clients, ok := spaceClients[spaceID]; ok {
		for _, clientChan := range clients {
			clientChan <- spaceData
		}
	}

	// Also broadcast to the general "spaces" channel for index page updates
	if clients, ok := spaceClients["spaces"]; ok {
		for _, clientChan := range clients {
			clientChan <- spaceData // You might want to send a different message structure here for the index page
		}
	}
}

func handleNewSpace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load spaces from file
	spaces, err := loadSpaces()
	if err != nil {
		http.Error(w, "Error loading spaces", http.StatusInternalServerError)
		return
	}

	newSpace := Space{
		ID:        uuid.New().String(),
		Name:      "New Space",
		CreatedAt: time.Now(),
	}

	spaces = append(spaces, newSpace)

	if err := saveSpaces(spaces); err != nil {
		http.Error(w, "Error saving spaces", http.StatusInternalServerError)
		return
	}

	// Generate and save password
	password := generatePassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		return
	}

	passwords, err := loadPasswords()
	if err != nil {
		http.Error(w, "Error loading passwords", http.StatusInternalServerError)
		return
	}

	passwords[newSpace.ID] = string(hashedPassword)

	if err := savePasswords(passwords); err != nil {
		http.Error(w, "Error saving passwords", http.StatusInternalServerError)
		return
	}

	spacesData, err := json.Marshal(spaces)
	if err != nil {
		http.Error(w, "Error marshaling spaces", http.StatusInternalServerError)
		return
	}

	broadcastSpaceUpdate("spaces", spacesData)

	http.Redirect(w, r, fmt.Sprintf("/spaces/%s?password=%s", newSpace.ID, password), http.StatusSeeOther)
}

func handleSpace(w http.ResponseWriter, r *http.Request) {
	spaceID := strings.TrimPrefix(r.URL.Path, "/spaces/")

	spaces, err := loadSpaces()
	if err != nil {
		http.Error(w, "Error loading spaces", http.StatusInternalServerError)
		return
	}

	var currentSpace *Space
	for _, space := range spaces {
		if space.ID == spaceID {
			currentSpace = &space
			break
		}
	}

	if currentSpace == nil {
		http.Error(w, "Space not found", http.StatusNotFound)
		return
	}

	// Serve the space.html template
	tmpl, err := template.ParseFiles("./client/space.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// You can pass data to the template if needed (e.g., the space name)
	err = tmpl.Execute(w, currentSpace)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func generatePassword() string {
	// For now, just a simple placeholder.  Improve later.
	return "password123"
}

func loadSpaces() ([]Space, error) {
	data, err := os.ReadFile("spaces.json")
	if err != nil {
		return nil, err
	}

	var spaces []Space
	if err := json.Unmarshal(data, &spaces); err != nil {
		return nil, err
	}
	return spaces, nil
}

func saveSpaces(spaces []Space) error {
	data, err := json.MarshalIndent(spaces, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("spaces.json", data, 0644)
}

func loadPasswords() (map[string]string, error) {
	data, err := os.ReadFile("passwords.json")
	if err != nil {
		return nil, err
	}

	var passwords map[string]string
	if err := json.Unmarshal(data, &passwords); err != nil {
		return nil, err
	}
	return passwords, nil
}

func savePasswords(passwords map[string]string) error {
	data, err := json.MarshalIndent(passwords, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("passwords.json", data, 0644)
}
