package svc

import (
	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/cloud-mill/cloudmill-websocket/internal/models"
	util "github.com/cloud-mill/cloudmill-websocket/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

var ClientPool *models.ClientPool

func init() {
	ClientPool = models.NewClientPool()

	go func() {
		err := ClientPool.Start()
		if err != nil {
			logger.Logger.Panic("error starting client pool", zap.Error(err))
		}
	}()
}

func AcceptConnection(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client_id")

	if len(clientId) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wsConnection, err := UpgradeHTTPToWS(w, r)
	if err != nil {
		logger.Logger.Error("error accepting connection", zap.Error(err))
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
		return
	}

	// TODO: an api call to database to register this client & mark this client is active

	// create a new client
	client := models.NewClient(clientId, wsConnection, ClientPool)

	// register client into pool
	ClientPool.Register <- client

	// make client listen for new messages
	go client.Read()
}

func HandleSendMessageToClients(clientIds []string, message models.Message) {
	for _, clientId := range clientIds {
		ClientPool.SendMessageToClient(clientId, message) // concurrency safe
	}
}
