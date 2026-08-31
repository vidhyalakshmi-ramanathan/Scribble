package collaboration

import (
	"errors"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
)

type Client struct {
	ID         string
	DocumentID string
	Conn       *websocket.Conn
	Send       chan []byte

	service *CollaborationService
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error", err)
			break
		}

		var wsMessage entities.WSMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.sendError(http.StatusBadRequest, "Invalid message format")
			continue
		}
		if wsMessage.Type != "edit" {
			c.sendError(http.StatusBadRequest, "Invalid message type")
			continue
		}

		log.Printf("EDIT received: document=%s user=%s message=%v", c.DocumentID, c.ID, wsMessage.Data)

		doc, err := c.service.ApplyEdit(c.ID, c.DocumentID, wsMessage.Data)
		if err != nil {

			log.Printf("ApplyEdit ERROR: %v", err)

			if appErr, ok := errors.AsType[*response.AppError](err); ok {
				c.sendError(appErr.Code, appErr.Message)
			} else {
				c.sendError(http.StatusInternalServerError, "Something went wrong")
			}
			continue
		}
		log.Printf("EDIT applied successfully")

		broadcastMsg := entities.BroadcastMessage{
			Type: "edit",
			Data: entities.BroadcastData{
				DocumentID: doc.DocumentID,
				Version:    doc.Version,
				Content:    doc.Content,
			},
		}
		data, err := json.Marshal(broadcastMsg)
		if err != nil {
			continue
		}

		hub.broadcast <- Message{
			DocumentID: c.DocumentID,
			Data:       data,
		}
	}
}

func (c *Client) writePump() {
	for message := range c.Send {
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)
		if err != nil {
			log.Println("Write error", err)
			return
		}
	}
}

func (c *Client) sendError(code int, message string) {
	msg := entities.WSMessage{
		Type: "error",
		Error: &entities.WSError{
			Code:    code,
			Message: message,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	c.Send <- data
}
