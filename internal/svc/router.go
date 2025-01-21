package svc

import (
	"net/http"
	"strings"

	"github.com/cloud-mill/cloudmill-websocket/internal/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter(
	authMiddleware func(next http.Handler, secretKey interface{}, authConfig AuthConfig) http.Handler,
) *mux.Router {

	// CORS config
	c := cors.New(cors.Options{
		AllowedOrigins:   config.Config.AllowedOrigins,
		ExposedHeaders:   []string{config.Config.Auth.CsrfHeaderName},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r := mux.NewRouter().StrictSlash(true)
	r.Use(c.Handler)

	for _, route := range CloudmillWebsocketOpenRoutes {
		r.Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	authConfig := AuthConfig{
		JwtCookieName:  config.Config.Auth.JwtCookieName,
		CsrfCookieName: config.Config.Auth.CsrfCookieName,
		CsrfHeaderName: config.Config.Auth.CsrfHeaderName,
	}

	for _, route := range CloudmillWebsocketProtectedRoutes {
		r.Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(authMiddleware(route.HandlerFunc, config.Config.Auth.AuthMiddlewareSecretKey, authConfig))
	}

	return r
}
