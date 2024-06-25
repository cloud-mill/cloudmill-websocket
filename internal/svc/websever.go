package svc

import (
	"github.com/cloud-mill/cloudmill-websocket/internal/config"
	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/cloud-mill/cm-common-golang/server"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func StartCloudmillWebsocket() {
	r := NewRouter(server.AuthMiddleware, server.ApiKeyMiddleware)
	http.Handle("/", r)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Config.Port), nil)

	if err != nil {
		logger.Logger.Panic("error starting web server", zap.Error(err))
	}

	logger.Logger.Info("CloudmillWebsocket is ready.")
}
