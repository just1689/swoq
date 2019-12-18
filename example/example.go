package main

import (
	"encoding/json"
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/swoq"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

var incQName = "incoming"

func main() {
	swoq.StartQueueClient()

	go createExampleIO()
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	swoq.StartWebServer(":8080", "/ws", incQName)
}

type Message struct {
	Name string `json:"name"`
}

func createExampleIO() {
	queue.Subscribe(incQName, func(m *nats.Msg) {
		w := &ws.WrappedMessage{}
		err := json.Unmarshal(m.Data, w)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		msg := &Message{}
		err = json.Unmarshal(w.Body, msg)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		reply := "hello " + msg.Name + "!"
		queue.Publish("client."+w.ClientID, reply)
	})
}
