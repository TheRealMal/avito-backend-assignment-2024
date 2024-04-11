package test

import (
	"avito-backend/internal/db"
	"avito-backend/internal/handlers"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUserBanner(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	s := handlers.NewServiceHandler(db.InitDB(context.Background(), mock), logger)
	token := fmt.Sprintf("Token %s", *handlers.NewJWT(handlers.UserRole))
	ts := httptest.NewServer(
		s.UserMiddleware(
			http.HandlerFunc(s.HandleUserBanner),
		),
	)

	t.Run("GetUserBanner", func(t *testing.T) {
		tagID, featID := 1, 1
		bannerContent, bannerIsActive := []byte("{}"), true
		mock.ExpectQuery("SELECT").
			WithArgs(tagID, featID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"b.content", "b.is_active"},
				).AddRow(
					bannerContent, bannerIsActive,
				),
			)
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d", ts.URL, tagID, featID),
			bytes.NewBuffer([]byte{}),
		)
		assert.Empty(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		res, err := client.Do(req)
		assert.Empty(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("GetCachedUserBanner", func(t *testing.T) {
		tagID, featID := 1, 1
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d", ts.URL, tagID, featID),
			bytes.NewBuffer([]byte{}),
		)
		assert.Empty(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		res, err := client.Do(req)
		assert.Empty(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("GetUnknownUserBanner", func(t *testing.T) {
		tagID, featID := 0, 0
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d", ts.URL, tagID, featID),
			bytes.NewBuffer([]byte{}),
		)
		assert.Empty(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		res, err := client.Do(req)
		assert.Empty(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("ForbiddenRole", func(t *testing.T) {
		token := fmt.Sprintf("Token %s", *handlers.NewJWT("FAKE_USER"))

		tagID, featID := 0, 0
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d", ts.URL, tagID, featID),
			bytes.NewBuffer([]byte{}),
		)
		assert.Empty(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		res, err := client.Do(req)
		assert.Empty(t, err)
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		token := ""

		tagID, featID := 0, 0
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d", ts.URL, tagID, featID),
			bytes.NewBuffer([]byte{}),
		)
		assert.Empty(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		res, err := client.Do(req)
		assert.Empty(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
}
