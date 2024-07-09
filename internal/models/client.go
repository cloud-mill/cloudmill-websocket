package models

import (
	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Id         string
	Conn       *websocket.Conn
	ClientPool *ClientPool
	writeMu    sync.Mutex
	readMu     sync.Mutex
}

func (c *Client) Read() {
	defer func() {
		c.ClientPool.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			logger.Logger.Info(
				"failed to close client connection",
				zap.String("clientId", c.Id),
				zap.Error(err),
			)
		}
	}()

	for {
		c.readMu.Lock()
		messageType, message, err := c.Conn.ReadMessage()
		c.readMu.Unlock()
		if err != nil {
			logger.Logger.Debug(
				"client connection error reading message",
				zap.String("clientId", c.Id),
				zap.Error(err),
			)
			return
		}

		err = c.HandleMessage(messageType, message)
		if err != nil {
			logger.Logger.Error(
				"error handling message",
				zap.String("clientId", c.Id),
				zap.Error(err),
			)
		}
	}
}

func (c *Client) Write(message ProcessedMessage) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if err := c.Conn.WriteJSON(message); err != nil {
		logger.Logger.Error(
			"failed to write message to client",
			zap.String("clientId", c.Id),
			zap.Error(err),
		)
	}
}

func (c *Client) Leave() {
	c.ClientPool.Unregister <- c
}

func NewClient(
	clientId string,
	conn *websocket.Conn,
	ClientPool *ClientPool,
) *Client {
	logger.Logger.Info("creating client", zap.String("clientId", clientId))
	return &Client{
		Id:         clientId,
		Conn:       conn,
		ClientPool: ClientPool,
	}
}

func (c *Client) HandleMessage(messageType int, message []byte) error {
	logger.Logger.Info(
		"client received message",
		zap.Int("type", messageType),
		zap.ByteString("message", message),
	)

	// TODO: do stuff about received message

	return nil
}
