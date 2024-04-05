package db

import "time"

type Banner struct {
	ID        int
	Tags      []int
	Feature   int
	Content   string
	IsActive  bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
