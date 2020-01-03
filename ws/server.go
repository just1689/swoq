package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	id := uuid.New().String()
	logrus.Infoln("New client: ", id)

	client := &Client{hub: hub, ClientID: id, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	client.unSub = hub.replier(client)

	go client.writePump()
	go client.readPump()
}
