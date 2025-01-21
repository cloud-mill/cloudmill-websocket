package svc

import (
	"net/http"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/cloud-mill/cloudmill-websocket/internal/models"
	"github.com/olahol/melody"
	"go.uber.org/zap"
)

var (
	ClientPool *models.ClientPool
	m          = melody.New()
)

func init() {
	ClientPool = models.NewClientPool()
	go func() {
		if err := ClientPool.Start(); err != nil {
			logger.Logger.Panic("error starting client pool", zap.Error(err))
		}
	}()

	m.HandleConnect(func(s *melody.Session) {
		clientId := s.Request.URL.Query().Get("client_id")
		if clientId == "" {
			logger.Logger.Warn("invalid client ID")

			if err := s.Close(); err != nil {
				logger.Logger.Error(
					"error closing session",
					zap.Error(err),
				)
			}
			return
		}

		client := &models.Client{
			Id:         clientId,
			Session:    s,
			ClientPool: ClientPool,
		}
		ClientPool.Register <- client

		logger.Logger.Info("client connected", zap.String("client_id", clientId))
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		clientId, exists := s.Get("client_id")
		if !exists {
			logger.Logger.Warn("message received from a session without client_id")
			return
		}

		client := ClientPool.GetClient(clientId.(string))
		if client == nil {
			logger.Logger.Warn(
				"client not found in the pool",
				zap.String("client_id", clientId.(string)),
			)
			return
		}

		if err := client.HandleMessage(msg); err != nil {
			logger.Logger.Error(
				"error handling client message",
				zap.String("client_id", clientId.(string)),
				zap.Error(err),
			)
		}
	})

	m.HandleDisconnect(func(s *melody.Session) {
		clientId, exists := s.Get("client_id")
		if !exists {
			logger.Logger.Warn("session disconnected without client_id")
			return
		}

		client := ClientPool.GetClient(clientId.(string))
		if client != nil {
			client.Leave()
		}

		logger.Logger.Info("client disconnected", zap.String("client_id", clientId.(string)))
	})
}

func AcceptConnection(w http.ResponseWriter, r *http.Request) {
	err := m.HandleRequest(w, r)
	if err != nil {
		logger.Logger.Error("error accepting websocket connection", zap.Error(err))
		http.Error(w, "error establishing connection", http.StatusInternalServerError)
	}
}
