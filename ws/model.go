package ws

import "encoding/json"

type WrappedMessage struct {
	ClientID string          `json:"clientID"`
	Body     json.RawMessage `json:"body"`
}
