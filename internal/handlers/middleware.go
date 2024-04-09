package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	MinAuthHeadersLen = 6 // "Token ..."
	AdminRole         = "ADMIN"
	UserRole          = "USER"
)

type ContextKey int

const (
	ContextRoleKey ContextKey = iota
)

var (
	jwtSecret []byte   = []byte("b112ebdd305b9191fe782f798fd922bf36784296a93720b8cc32469f263bc670")
	roles     []string = []string{AdminRole, UserRole}
)

type UniqueClaims struct {
	jwt.StandardClaims
	TokenId string `json:"jti,omitempty"`
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := []rune(r.Header.Get("Authorization"))
		if len(token) <= MinAuthHeadersLen {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		role, err := checkToken(string(token[MinAuthHeadersLen:]), AdminRole)
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

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := []rune(r.Header.Get("Authorization"))
		if len(token) <= MinAuthHeadersLen {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		role, err := checkToken(string(token[MinAuthHeadersLen:]), AdminRole)
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

func checkToken(inputToken string, checkRole string) (string, error) {
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
		fmt.Println(err)
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
	// Get field "role" and verify
	userRole, ok := payload["role"].(string)
	if !ok {
		return "", errors.New("bad token")
	}
	return userRole, nil
}

func NewJWT(role string) *string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(60 * time.Minute),
		"role": role,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil
	}
	return &tokenString
}
