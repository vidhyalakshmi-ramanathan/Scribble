package collaboration

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/middleware"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigin := strings.TrimRight(
			os.Getenv("FRONTEND_ORIGIN"),
			"/",
		)

		// Local development fallback. Configure FRONTEND_ORIGIN
		// to your deployed frontend's origin in production.
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000"
		}

		return origin != "" && origin == allowedOrigin
	},
}

type Handler struct {
	roleService          services.RoleService
	collaborationService *CollaborationService
	hub                  *Hub
}

func NewCollaborationHandler(roleService services.RoleService, collaborationService *CollaborationService, hub *Hub) *Handler {
	return &Handler{roleService: roleService, collaborationService: collaborationService, hub: hub}
}

func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// authentication + document access checks
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	_, err := h.roleService.CheckMember(claims.UserID, docID)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	//Upgrade HTTP to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}

	// defer ws.Close()

	client := &Client{
		ID:         claims.UserID,
		DocumentID: docID,
		Conn:       ws,
		Send:       make(chan []byte, 256),
		service:    h.collaborationService,
	}
	h.hub.register <- client

	go client.writePump(h.hub)
	client.readPump(h.hub)
}

// 	// Temporary echo test
// 	defer ws.Close()
// 	for {
// 		messageType, payload, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			return
// 		}
// 		// Echo the same message back to the client
// 		err = ws.WriteMessage(messageType, payload)
// 		if err != nil {
// 			log.Println("Write error:", err)
// 			return
// 		}
// 	}
// }
