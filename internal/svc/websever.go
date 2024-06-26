package svc

import (
	"github.com/cloud-mill/cloudmill-websocket/internal/config"
	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func StartCloudmillWebsocket() {
	r := NewRouter(AuthMiddleware)
	http.Handle("/", r)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Config.Port), nil)

	if err != nil {
		logger.Logger.Panic("error starting web server", zap.Error(err))
	}

	logger.Logger.Info("CloudmillWebsocket is ready.")
}
