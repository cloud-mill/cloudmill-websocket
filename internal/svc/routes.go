package svc

import (
	"net/http"
)

type Route struct {
	Name         string
	Description  string
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	Authenticate bool
}

type Routes []Route

var CloudmillWebsocketProtectedRoutes = Routes{
	Route{
		Name:        "connect",
		Method:      "GET,OPTIONS",
		Pattern:     "/connect",
		HandlerFunc: AcceptConnection,
	},
}

var okHandler = func(writer http.ResponseWriter, r *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		return
	}
}

var CloudmillWebsocketOpenRoutes = Routes{
	Route{
		Name:        "HealthCheck",
		Method:      "GET,OPTIONS",
		Pattern:     "/healthz",
		HandlerFunc: okHandler,
	},

	Route{
		Name:        "ReadinessCheck",
		Method:      "GET,OPTIONS",
		Pattern:     "/readyz",
		HandlerFunc: okHandler,
	},

	Route{
		Name:        "LivenessCheck",
		Method:      "GET,OPTIONS",
		Pattern:     "/livez",
		HandlerFunc: okHandler,
	},
}
