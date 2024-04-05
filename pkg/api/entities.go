package api

import (
	"avito-backend/pkg/db"
)

type Service struct {
	db    db.Database
	cache interface{}
}
