package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

const (
	MinAuthHeadersLen = 6 // "Token ..."
	AdminRole         = "ADMIN"
	UserRole          = "USER"
)

var jwtSecret = []byte("b112ebdd305b9191fe782f798fd922bf36784296a93720b8cc32469f263bc670")

type ContextKey int

const ContextRoleKey ContextKey = iota

func (s ServiceHandler) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Info("req_details",
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}

func (s ServiceHandler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := []rune(r.Header.Get("Authorization"))
		if len(token) <= MinAuthHeadersLen {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		role, err := checkToken(string(token[MinAuthHeadersLen:]))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if role != AdminRole {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s ServiceHandler) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := []rune(r.Header.Get("Authorization"))
		if len(token) <= MinAuthHeadersLen {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		role, err := checkToken(string(token[MinAuthHeadersLen:]))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if role != AdminRole && role != UserRole {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextRoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkToken(inputToken string) (string, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, errors.New("bad sign method")
		}
		return jwtSecret, nil
	}

	// Parse incoming string
	token, err := jwt.Parse(inputToken, hashSecretGetter)
	if err != nil {
		return "", errors.New("bad token")
	}
	if !token.Valid {
		return "", nil
	}

	// Unpack token claims
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("bad token")
	}
	// Get field role and expiration time
	userRole, ok := payload["role"].(string)
	if !ok {
		return "", errors.New("bad token")
	}
	return userRole, nil
}

func NewJWT(role string) *string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "0",
		"iat":  time.Now().Unix(),
		"role": role,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil
	}
	return &tokenString
}
