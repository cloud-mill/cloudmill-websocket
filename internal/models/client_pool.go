package models

import (
	"sync"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"go.uber.org/zap"
)

type ClientPool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	rwMutex    sync.RWMutex
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (cp *ClientPool) Start() error {
	for {
		select {
		case client := <-cp.Register:
			logger.Logger.Info("registering client", zap.String("client_id", client.Id))
			cp.SetClient(client.Id, client)

		case client := <-cp.Unregister:
			logger.Logger.Info("unregistering client", zap.String("client_id", client.Id))
			cp.DeleteClient(client.Id)
		}
	}
}

func (cp *ClientPool) GetClient(clientId string) *Client {
	cp.rwMutex.RLock()
	defer cp.rwMutex.RUnlock()

	return cp.Clients[clientId]
}

func (cp *ClientPool) SetClient(clientId string, client *Client) {
	cp.rwMutex.Lock()
	defer cp.rwMutex.Unlock()

	cp.Clients[clientId] = client
}

func (cp *ClientPool) DeleteClient(clientId string) {
	cp.rwMutex.Lock()
	defer cp.rwMutex.Unlock()

	delete(cp.Clients, clientId)
}

func (cp *ClientPool) SendMessageToClient(clientId string, message Message) {
	client := cp.GetClient(clientId)
	if client != nil {
		processedMessage := ProcessedMessage{
			Id:        message.Id,
			Type:      message.Type,
			Timestamp: message.Timestamp,
			Payload:   message.Payload,
		}
		client.Write(processedMessage)
	}
}

func (cp *ClientPool) ClientExitFromPool(clientId string) {
	client := cp.GetClient(clientId)
	if client != nil {
		cp.Unregister <- client

		if err := client.Session.Close(); err != nil {
			logger.Logger.Error(
				"failed to close client session",
				zap.String("client_id", client.Id),
				zap.Error(err),
			)
		}
	}
}
