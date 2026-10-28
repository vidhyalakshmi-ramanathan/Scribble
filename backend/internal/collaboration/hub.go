package collaboration

import "log"

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
		broadcast:  make(chan Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			select {
			case <-client.done:
				continue
			default:
			}

			if h.clients[client.DocumentID] == nil {
				h.clients[client.DocumentID] = make(map[*Client]bool)
			}
			h.clients[client.DocumentID][client] = true

		case client := <-h.unregister:
			h.removeClient(client)
			client.stop()

		case message := <-h.broadcast:
			for client := range h.clients[message.DocumentID] {
				select {
				case <-client.done:
					h.removeClient(client)
				case client.Send <- message.Data:
					// Message queued for the client's writePump.
				default:
					log.Printf("Client %s is too slow, disconnecting", client.ID)
					h.removeClient(client)
					client.stop()
				}
			}
		}
	}
}

func (h *Hub) removeClient(client *Client) bool {
	clients, ok := h.clients[client.DocumentID]
	if !ok {
		return false
	}

	if _, exists := clients[client]; !exists {
		return false
	}

	delete(clients, client)

	if len(clients) == 0 {
		delete(h.clients, client.DocumentID)
	}

	return true
}
