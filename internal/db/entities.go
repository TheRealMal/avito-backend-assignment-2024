package db

import (
	"time"

	"github.com/goccy/go-json"
)

type Banner struct {
	ID        int             `json:"id"`
	Tags      []int           `json:"tags"`
	Feature   int             `json:"feature"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}
