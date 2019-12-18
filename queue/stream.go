package queue

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var DefaultConn *nats.Conn

func BuildDefaultConn() {
	var err error
	DefaultConn, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Println("NATS connection established.")
}

func GetPublisher(subj string) func(body []byte) {
	return func(body []byte) {
		DefaultConn.Publish(subj, body)
	}
}

func Publish(subj, body string) {
	DefaultConn.Publish(subj, []byte(body))
}

func Subscribe(subj string, handler func(m *nats.Msg)) (close func()) {

	sub, err := DefaultConn.Subscribe(subj, handler)
	if err != nil {
		//TODO: figure out how we can handle this? Retry?
		logrus.Panicln(err)
	}
	close = func() {
		sub.Unsubscribe()
		sub.Drain()
	}
	return

}
