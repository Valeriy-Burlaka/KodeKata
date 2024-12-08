package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Space struct {
	ID      string   `json:"id"`
	Rows    int      `json:"rows"`
	Cols    int      `json:"cols"`
	Enabled []string `json:"enabled"`
}

type Store struct {
	mu         sync.RWMutex
	spaces     map[string]*Space
	spaceLocks map[string]*sync.RWMutex
	path       string
}

func NewStore(path string) (*Store, error) {
	s := &Store{
		mu:         sync.RWMutex{},
		spaces:     make(map[string]*Space),
		spaceLocks: make(map[string]*sync.RWMutex),
		path:       path,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for id := range s.spaces {
		s.spaceLocks[id] = &sync.RWMutex{}
	}

	go s.periodicallyDump()

	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.spaces)
}

func (s *Store) periodicallyDump() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		if err := s.dump(); err != nil {
			log.Printf("Failed to dump spaces: %v", err)
		}
	}
}

func (s *Store) dump() error {
	s.mu.RLock()
	data, err := json.Marshal(s.spaces)
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) SaveSpace(space *Space) error {
	s.spaceLocks[space.ID] = &sync.RWMutex{}
	s.spaces[space.ID] = space
	return nil
}

func (s *Store) GetSpace(id string) (*Space, error) {
	lock, ok := s.spaceLocks[id]
	if !ok {
		return nil, fmt.Errorf("space not found: %s", id)
	}

	lock.RLock()
	space := s.spaces[id]
	lock.RUnlock()

	return space, nil
}
