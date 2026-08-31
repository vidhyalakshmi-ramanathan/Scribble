package dto

import (
	"github.com/golang-jwt/jwt/v5"
	auth "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
)

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ClaimsData struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type DashboardResponse struct {
	MyDocuments     []doc.Documents        `json:"my_documents"`
	RecentDocuments []doc.Documents        `json:"recent_documents"`
	SharedDocuments []doc.Documents        `json:"shared_documents"`
	PendingInvites  []auth.DocumentInvites `json:"pending_invites"`
}
