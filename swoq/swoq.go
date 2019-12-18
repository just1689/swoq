package swoq

import (
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

var hub *ws.Hub

func StartQueueClient(incomingMsgQueueName string) {
	hub = ws.NewHub(queue.GetPublisher(incomingMsgQueueName), StartReplier)
	go hub.Run()
}

func StartWebServer(listenAddr string, wsUrl string) {
	queue.BuildDefaultConn()

	ConfigureWebServer(wsUrl)
	logrus.Panic(http.ListenAndServe(listenAddr, nil))
}

func ConfigureWebServer(wsUrl string) {
	http.HandleFunc("/", ws.ServeHome)
	http.HandleFunc(wsUrl, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
}

//StartReplier reads from the Queue and writes to the hub
func StartReplier(id string) {
	found, client := hub.GetClientByID(id)
	if !found {
		//TODO: Should exist at this time?
		logrus.Errorln("failed to start replier - client not found by id: ", id)
		return
	}
	queue.Subscribe(id, func(m *nats.Msg) {
		client.Send(m.Data)
	})

}
