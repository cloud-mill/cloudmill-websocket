package svc

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtClaim struct {
	CustomClaims UserCustomClaim
	jwt.RegisteredClaims
}

type UserCustomClaim struct {
	AccountId uuid.UUID
	Username  string
	Email     string
}

type AuthConfig struct {
	JwtCookieName  string
	CsrfCookieName string
	CsrfHeaderName string
}

func ConvertToByteSecretKey(secretKey interface{}) ([]byte, error) {
	switch secretKey := secretKey.(type) {
	case string:
		return []byte(secretKey), nil
	case []byte:
		return secretKey, nil
	default:
		return nil, fmt.Errorf("invalid secret key type")
	}
}

func validateJWTAndGetClaims(
	token string,
	claims *JwtClaim,
	secretKey interface{},
) (jwt.Claims, error) {
	byteSecretKey, err := ConvertToByteSecretKey(secretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key type")
	}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return byteSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	return tkn.Claims, nil
}

func AuthMiddleware(
	next http.Handler,
	secretKey interface{},
	authConfig AuthConfig,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authConfig.JwtCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// validate JWT Cookie
		token := cookie.Value
		jwtClaims := JwtClaim{}
		_, err = validateJWTAndGetClaims(
			token,
			&jwtClaims,
			secretKey,
		)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get CSRF cookie
		csrfCookie, err := r.Cookie(authConfig.CsrfCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// validate CSRF
		csrfId := r.Header.Get(authConfig.CsrfHeaderName)
		if csrfId == "" || csrfId != csrfCookie.Value {
			log.Print("WARN: CSRF ATTACK")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
