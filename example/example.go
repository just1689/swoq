package main

import (
	"encoding/json"
	"flag"
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/swoq"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var incQName = "incoming"

func main() {
	flag.Parse()
	swoq.StartQueueClient()

	go createExampleIO()
	swoq.StartWebServer(":8080", "/ws", incQName)
}

func createExampleIO() {
	queue.Subscribe(incQName, func(m *nats.Msg) {
		w := &ws.WrappedMessage{}
		err := json.Unmarshal(m.Data, w)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		reply := "hello " + string(w.Body) + "!"
		queue.Publish(w.ClientID, reply)
	})
}
