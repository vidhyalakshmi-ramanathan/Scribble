package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
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
	MyDocuments     []MyDocuments     `json:"myDocuments"`
	RecentDocuments []RecentDocuments `json:"recentDocuments"`
	SharedDocuments []SharedDocuments `json:"sharedDocuments"`
	PendingInvites  []PendingInvites  `json:"pendingInvites"`
}

type MyDocuments struct {
	DocumentID string    `json:"documentID"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type RecentDocuments struct {
	DocumentID string    `json:"documentID"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	UpdatedAt  time.Time `json:"updatedAt"`
	OwnerID    string    `json:"ownerID"`
}

type SharedDocuments struct {
	DocumentID string      `json:"documentID"`
	Title      string      `json:"title"`
	Category   string      `json:"category"`
	UpdatedAt  time.Time   `json:"updatedAt"`
	SharedBy   string      `json:"sharedBy"`
	SharedByID string      `json:"sharedByID"`
	Role       enums.Roles `json:"role"`
}

type PendingInvites struct {
	InvitationID string        `json:"invitationID"`
	DocumentID   string        `json:"documentID"`
	Title        string        `json:"title"`
	Role         enums.Roles   `json:"role"`
	InvitedEmail string        `json:"invitedEmail"`
	InviterName  string        `json:"inviteName"`
	InvitedBy    string        `json:"invitedBy"`
	Status       enums.Invites `json:"status"`
	ExpiresAt    time.Time     `json:"expiresAt"`
	CreatedAt    time.Time     `json:"createdAt"`
}
