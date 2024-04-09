package handlers

import (
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
		authorized, err := checkToken(string(token[MinAuthHeadersLen:]), AdminRole)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !authorized {
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
		var err error
		var authorized bool
		for _, role := range roles {
			authorized, err = checkToken(string(token[MinAuthHeadersLen:]), role)
			if err != nil {
				continue
			}
			if !authorized {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})
}

func checkToken(inputToken string, checkRole string) (bool, error) {
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
		return false, errors.New("bad token")
	}
	if !token.Valid {
		return false, nil
	}

	// Unpack token claims
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("bad token")
	}
	// Get field "role" and verify
	userRole, ok := payload["role"].(string)
	if !ok {
		return false, errors.New("bad token")
	}
	if userRole != checkRole {
		return false, nil
	}
	return true, nil
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
