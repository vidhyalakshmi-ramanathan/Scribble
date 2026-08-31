package collaboration

type Hub struct {
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

type Message struct {
	DocumentID string
	Data       []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if h.clients[client.DocumentID] == nil {
				h.clients[client.DocumentID] = make(map[*Client]bool)
			}
			h.clients[client.DocumentID][client] = true

		case client := <-h.unregister:
			if clients, ok := h.clients[client.DocumentID]; ok {
				delete(clients, client)

				if len(clients) == 0 {
					delete(h.clients, client.DocumentID)
				}
			}
			close(client.Send)

		case message := <-h.broadcast:
			for client := range h.clients[message.DocumentID] {
				client.Send <- message.Data
			}
		}
	}
}
