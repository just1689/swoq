package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MessageOut struct {
	Name string `json:"name"`
}

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		time.Sleep(1 * time.Second)
		c.WriteJSON(MessageOut{Name: "Justin"})
	}()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	go func() {
		for {
			time.Sleep(1500 * time.Millisecond)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter id to inherit: ")
			text, _ := reader.ReadString('\n')
			text = strings.ReplaceAll(text, "\n", "")
			msg := ".id." + text
			fmt.Println(text)
			c.WriteMessage(websocket.BinaryMessage, []byte(msg))
		}
	}()

	select {}
}
