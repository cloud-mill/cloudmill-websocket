package svc

import (
	"github.com/cloud-mill/cloudmill-websocket/internal/config"
	"net/http"
	"strings"

	"github.com/cloud-mill/cm-common-golang/models"
	pb "github.com/cloud-mill/cm-common-golang/proto/protos/out"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type AuthMiddleware func(next http.Handler, secretKey interface{}, authConfig models.AuthConfig) http.Handler
type ApiKeyMiddleware func(next http.Handler, dbServiceClient pb.MissionHandlerClient) http.Handler

func NewRouter(authMiddleware AuthMiddleware, apiMiddleware ApiKeyMiddleware) *mux.Router {

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

	authConfig := models.AuthConfig{
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
