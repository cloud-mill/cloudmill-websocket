package svc

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	readBufferBytesSize  = 4096
	writeBufferBytesSize = 4096
)

func makeUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  readBufferBytesSize,
		WriteBufferSize: writeBufferBytesSize,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
}

func UpgradeHTTPToWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := makeUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, err
}
