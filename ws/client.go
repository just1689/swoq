package ws

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
)

type Client struct {
	hub      *Hub
	ClientID string
	conn     *websocket.Conn
	send     chan []byte
	unSub    func()
}

func (c *Client) Send(b []byte) {
	c.send <- b
}

func (c *Client) GetReplyQueueName() string {
	return "client." + c.ClientID
}

func (c *Client) RemoveMeCleanly() {
	logrus.Println("Removing client ", c.ClientID)
	c.hub.unregister <- c
	c.conn.Close()
	c.unSub()
}

func (c *Client) readPump() {
	defer func() {
		c.RemoveMeCleanly()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Println("Could not read as the socket is closed")
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if string(message[0:1]) == "." {
			commands := strings.Split(string(message), ".")
			c.handleCommands(commands[1], commands[2])
			continue
		}

		var b []byte
		if b, err = json.Marshal(WrappedMessage{
			ClientID: c.ClientID,
			Body:     message,
		}); err != nil {
			logrus.Errorln(err)
			continue
		}
		c.hub.publisher(b)

	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.RemoveMeCleanly()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleCommands(key, value string) {
	if key == "id" {
		c.unSub()
		logrus.Println("Client reconnected: ", value)
		c.ClientID = value
		c.unSub = c.hub.replier(c)
	}

}
