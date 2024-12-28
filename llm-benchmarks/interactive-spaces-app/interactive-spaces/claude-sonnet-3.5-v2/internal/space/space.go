package space

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/maps"
)

const (
	GridSize      = 100
	PasswordBytes = 16
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Space struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	CreatedAt       time.Time       `json:"created_at"`
	PasswordHash    string          `json:"password_hash"` // bcrypt hash
	Started         bool            `json:"started"`
	CurrentPosition Position        `json:"current_position"`
	ActiveCells     map[string]bool `json:"active_cells"` // "x,y" -> true
	ConnectedCount  int             `json:"connected_count"`
	mu              sync.RWMutex    // protects all fields
}

// NewSpace creates a new space with generated ID and password
func NewSpace(name string) (*Space, string, error) {
	// Generate random ID (16 bytes, base64 encoded)
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return nil, "", err
	}
	spaceID := base64.URLEncoding.EncodeToString(id)

	// Generate random password (16 bytes, base64 encoded)
	pwd := make([]byte, PasswordBytes)
	if _, err := rand.Read(pwd); err != nil {
		return nil, "", err
	}
	password := base64.URLEncoding.EncodeToString(pwd)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	space := &Space{
		ID:              spaceID,
		Name:            name,
		CreatedAt:       time.Now(),
		PasswordHash:    string(hash),
		Started:         false,
		CurrentPosition: Position{0, 0},
		ActiveCells:     make(map[string]bool),
		ConnectedCount:  0,
	}

	return space, password, nil
}

// cellKey generates a string key for the ActiveCells map
func cellKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// Evolve adds one new active cell at the current position
func (s *Space) Evolve() (Position, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Started {
		return s.CurrentPosition, false
	}

	// Mark current cell as active
	key := cellKey(s.CurrentPosition.X, s.CurrentPosition.Y)
	s.ActiveCells[key] = true

	// Move to next position
	s.CurrentPosition.X++
	if s.CurrentPosition.X >= GridSize {
		s.CurrentPosition.X = 0
		s.CurrentPosition.Y++
		if s.CurrentPosition.Y >= GridSize {
			s.CurrentPosition.Y = 0 // Wrap around
		}
	}

	return s.CurrentPosition, true
}

// Start begins space evolution
func (s *Space) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Started = true
}

// Stop pauses space evolution
func (s *Space) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Started = false
}

// Reset clears all active cells and resets position
func (s *Space) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActiveCells = make(map[string]bool)
	s.CurrentPosition = Position{0, 0}
}

// UpdateName updates space name if the password is correct
func (s *Space) UpdateName(name, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := bcrypt.CompareHashAndPassword([]byte(s.PasswordHash), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	s.Name = name
	return nil
}

// ConnectionOpened increments the connected client count
func (s *Space) ConnectionOpened() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ConnectedCount++
}

// ConnectionClosed decrements the connected client count
func (s *Space) ConnectionClosed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ConnectedCount > 0 {
		s.ConnectedCount--
	}
}

// GetState returns a copy of space state for reading
func (s *Space) GetState() SpaceState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return SpaceState{
		ID:              s.ID,
		Name:            s.Name,
		CreatedAt:       s.CreatedAt,
		Started:         s.Started,
		CurrentPosition: s.CurrentPosition,
		ActiveCells:     maps.Clone(s.ActiveCells),
		ConnectedCount:  s.ConnectedCount,
	}
}
