package swoq

import (
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"log"
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
	http.HandleFunc("/", ServeHome)
	http.HandleFunc(wsUrl, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
}

//StartReplier reads from the Queue and writes to the hub
func StartReplier(client *ws.Client) {
	queue.Subscribe(client.GetReplyQueueName(), func(m *nats.Msg) {
		client.Send(m.Data)
	})

}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
