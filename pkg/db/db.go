package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	QueryGetUserBanner             = "SELECT b.content FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id=$2 AND b.feature=$1;"
	QueryGetBanners                = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id GROUP BY b.id LIMIT $1 OFFSET $2;"
	QueryGetBannersFilterByFeature = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id GROUP BY b.id HAVING b.feature = $1 LIMIT $2 OFFSET $3;"
	QueryGetBannersFilterByTag     = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id = $1 GROUP BY b.id LIMIT $2 OFFSET $3;"
	QueryGetBannersFilterByBoth    = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id = $1 GROUP BY b.id HAVING b.feature = $2 LIMIT $3 OFFSET $4;"
)

type SQLDatabase struct {
	pool *pgxpool.Pool
}

func InitDB(databaseURL string) (SQLDatabase, error) {
	res := SQLDatabase{}
	var err error

	res.pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return res, err
	}

	return res, nil
}

type Database interface {
	GetBannerContent(tagID int, featureID int) (*[]byte, error)
	GetBanners(tagID int, featureID int, limit int, offset int) ([]Banner, error)
	CreateBanner(banner Banner) error
	UpdateBanner(id int, tagIDs []int, featureID int, content []byte, isActive bool) error
	DeleteBanner(id int) error
}

func (db SQLDatabase) GetBannerContent(tagID int, featureID int) (*[]byte, error) {
	var content []byte
	err := db.pool.QueryRow(
		context.Background(),
		QueryGetUserBanner,
		tagID,
		featureID,
	).Scan(&content)

	switch err {
	case nil:
		return &content, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (db SQLDatabase) GetBanners(tagID int, featureID int, limit int, offset int) ([]Banner, error) {
	result := []Banner{}
	rows, err := db.chooseQuery(tagID, featureID, limit, offset)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var banner Banner
		err := rows.Scan(
			&banner.ID,
			&banner.Feature,
			&banner.Content,
			&banner.CreatedAt,
			&banner.UpdatedAt,
			&banner.IsActive,
			&banner.Tags,
		)
		if err != nil {
			return result, err
		}
		result = append(result, banner)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func (db SQLDatabase) chooseQuery(tagID int, featureID int, limit int, offset int) (pgx.Rows, error) {
	switch {
	case tagID != -1 && featureID != -1:
		return db.pool.Query(
			context.Background(),
			QueryGetBannersFilterByBoth,
			tagID,
			featureID,
			limit,
			offset,
		)
	case tagID != -1:
		return db.pool.Query(
			context.Background(),
			QueryGetBannersFilterByTag,
			tagID,
			limit,
			offset,
		)
	case featureID != -1:
		return db.pool.Query(
			context.Background(),
			QueryGetBannersFilterByFeature,
			featureID,
			limit,
			offset,
		)
	default:
		return db.pool.Query(
			context.Background(),
			QueryGetBanners,
			limit,
			offset,
		)
	}
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
