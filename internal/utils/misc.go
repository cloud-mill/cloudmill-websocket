package utils

import (
	"net/http"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"go.uber.org/zap"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		logger.Logger.Error("error writing json response", zap.Error(err))
	}
}
