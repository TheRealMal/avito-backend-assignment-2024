package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	QueryGetUserBanner             = "SELECT b.content FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id=$1 AND b.feature=$2;"
	QueryGetBanners                = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id GROUP BY b.id LIMIT $1 OFFSET $2;"
	QueryGetBannersFilterByFeature = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id GROUP BY b.id HAVING b.feature = $1 LIMIT $2 OFFSET $3;"
	QueryGetBannersFilterByTag     = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id = $1 GROUP BY b.id LIMIT $2 OFFSET $3;"
	QueryGetBannersFilterByBoth    = "SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id WHERE t.id = $1 GROUP BY b.id HAVING b.feature = $2 LIMIT $3 OFFSET $4;"
	QueryInsertBanner              = "INSERT INTO banners (is_active, feature, content) VALUES ($1, $2, convert_to($3, 'LATIN1')) RETURNING id;"
	QueryInsertBannerTag           = "INSERT INTO tags (id, banner_id) VALUES ($1, $2);"
	QueryDeleteBannerTags          = "DELETE FROM tags WHERE banner_id = $1;"
	QueryDeleteBanner              = "DELETE FROM banners WHERE id = $1;"
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
	UpdateBanner(id int, tagIDs *[]int, featureID *int, content *json.RawMessage, isActive *bool) (bool, error)
	DeleteBanner(id int) (bool, error)
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
	var id int
	err := db.pool.QueryRow(
		context.Background(),
		QueryInsertBanner,
		banner.IsActive,
		banner.Feature,
		banner.Content,
	).Scan(&id)
	if err != nil {
		return err
	}
	for _, tag := range banner.Tags {
		_, err := db.pool.Exec(
			context.Background(),
			QueryInsertBannerTag,
			tag,
			id,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db SQLDatabase) UpdateBanner(id int, tagIDs *[]int, featureID *int, content *json.RawMessage, isActive *bool) (bool, error) {
	// Prepare parametrized query and update not null columns
	if featureID != nil || content != nil || isActive != nil {
		query, args := prepareQuery(id, featureID, content, isActive)
		res, err := db.pool.Exec(
			context.Background(),
			query,
			args...,
		)
		if err != nil {
			return false, err
		}
		if res.RowsAffected() == 0 {
			return false, nil
		}
	}

	if tagIDs == nil {
		return true, nil
	}
	// Delete existing tags
	res, err := db.pool.Exec(
		context.Background(),
		QueryDeleteBannerTags,
		id,
	)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 0 {
		return false, nil
	}
	// Insert new tags
	for _, tag := range *tagIDs {
		_, err := db.pool.Exec(
			context.Background(),
			QueryInsertBannerTag,
			tag,
			id,
		)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func prepareQuery(id int, featureID *int, content *json.RawMessage, isActive *bool) (string, []any) {
	query, argsCounter := strings.Builder{}, 1
	args := make([]any, 0)
	queryParts := make([]string, 0)
	query.WriteString("UPDATE banners SET ")
	for i := 0; i < 3; i++ {
		switch {
		case i == 0 && featureID != nil:
			queryParts = append(queryParts, fmt.Sprintf("feature=$%d", argsCounter))
			args = append(args, *featureID)
			argsCounter++
		case i == 1 && content != nil:
			queryParts = append(queryParts, fmt.Sprintf("content=$%d", argsCounter))
			args = append(args, *content)
			argsCounter++
		case i == 2 && isActive != nil:
			queryParts = append(queryParts, fmt.Sprintf("is_active=$%d", argsCounter))
			args = append(args, *isActive)
			argsCounter++
		}
	}
	query.WriteString(strings.Join(queryParts, ", "))
	query.WriteString(fmt.Sprintf(" WHERE id=$%d", argsCounter))
	args = append(args, id)
	return query.String(), args
}

func (db SQLDatabase) DeleteBanner(id int) (bool, error) {
	// Delete existing tags
	res, err := db.pool.Exec(
		context.Background(),
		QueryDeleteBannerTags,
		id,
	)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 0 {
		return false, nil
	}

	res, err = db.pool.Exec(
		context.Background(),
		QueryDeleteBanner,
		id,
	)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}
