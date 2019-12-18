package swoq

import (
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

var hub *ws.Hub

func StartQueueClient() {
	queue.BuildDefaultConn()
}

func StartWebServer(listenAddr string, wsUrl string, incomingMsgQueueName string) {
	hub = ws.NewHub(queue.GetPublisher(incomingMsgQueueName), startReplier)
	go hub.Run()
	ConfigureWebServer(wsUrl)
	logrus.Panic(http.ListenAndServe(listenAddr, nil))
}

func ConfigureWebServer(wsUrl string) {
	http.HandleFunc(wsUrl, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
}

func startReplier(client *ws.Client) {
	queue.Subscribe(client.GetReplyQueueName(), func(m *nats.Msg) {
		client.Send(m.Data)
	})

}
