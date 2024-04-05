package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type SQLDatabase struct {
	C *pgx.Conn
}

func InitDB(databaseURL string) Database {
	res := SQLDatabase{}
	var err error

	res.C, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil
	}

	return res
}

type Database interface {
	GetUserBanners(tagID int, featureID int, useLastRevision bool) ([]byte, error)
	GetBanners(tagID int, featureID int, limit int, offset int) ([]Banner, error)
	CreateBanner(banner Banner) error
	UpdateBanner(id int, tagIDs []int, featureID int, content []byte, isActive bool) error
	DeleteBanner(id int) error
}

func (db SQLDatabase) GetUserBanners(tagID int, featureID int, useLastRevision bool) ([]byte, error) {
	return []byte{}, nil
}

func (db SQLDatabase) GetBanners(tagID int, featureID int, limit int, offset int) ([]Banner, error) {
	return []Banner{}, nil
}

func (db SQLDatabase) CreateBanner(banner Banner) error {
	return nil
}

func (db SQLDatabase) UpdateBanner(id int, tagIDs []int, featureID int, content []byte, isActive bool) error {
	return nil
}

func (db SQLDatabase) DeleteBanner(id int) error {
	return nil
}
