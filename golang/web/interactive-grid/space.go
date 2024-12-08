package main

import (
	"time"
)

type Space struct {
	ID        string    `json:"id"`
	Rows      int       `json:"rows"`
	Cols      int       `json:"cols"`
	Enabled   []string  `json:"enabled"`    // list of enabled cells
	CreatedAt time.Time `json:"created_at"` // time the space was created
}
