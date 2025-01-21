package models

import (
	"encoding/json"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/olahol/melody"
	"go.uber.org/zap"
)

type Client struct {
	Id         string
	Session    *melody.Session
	ClientPool *ClientPool
}

func NewClient(
	clientId string,
	session *melody.Session,
	clientPool *ClientPool,
) *Client {
	logger.Logger.Info("new client", zap.String("clientId", clientId))
	return &Client{
		Id:         clientId,
		Session:    session,
		ClientPool: clientPool,
	}
}

func (c *Client) Write(message ProcessedMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error(
			"failed to encode message to JSON",
			zap.String("client_id", c.Id),
			zap.Error(err),
		)
		return
	}

	if err := c.Session.Write(messageBytes); err != nil {
		logger.Logger.Error(
			"failed to write message to client",
			zap.String("client_id", c.Id),
			zap.Error(err),
		)
	}
}

func (c *Client) Leave() {
	c.ClientPool.Unregister <- c
	logger.Logger.Info("client left the pool", zap.String("client_id", c.Id))
}

func (c *Client) HandleMessage(message []byte) error {
	logger.Logger.Info(
		"client received message",
		zap.String("clientId", c.Id),
		zap.ByteString("message", message),
	)

	// TODO: do stuff about message
	return nil
}
