package svc

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/cloud-mill/cloudmill-websocket/internal/config"
	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"go.uber.org/zap"
)

func StartCloudmillWebsocket() {
	router := NewRouter(AuthMiddleware)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Config.Port),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logger.Logger.Info("starting server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Panic("error starting server", zap.Error(err))
		}
	}()

	<-stop
	logger.Logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Error("error during server shutdown", zap.Error(err))
	}
	logger.Logger.Info("server stopped")
}
