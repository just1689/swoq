package ws

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	publisher  func(body []byte)
	replier    func(client *Client) func()
}

func NewHub(publisher func(body []byte), replier func(client *Client) func()) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		publisher:  publisher,
		replier:    replier,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}
