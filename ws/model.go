package ws

import "encoding/json"

type WrappedMessage struct {
	ClientID string          `json:"clientID"`
	Title    string          `json:"title"`
	Body     json.RawMessage `json:"body"`
}
