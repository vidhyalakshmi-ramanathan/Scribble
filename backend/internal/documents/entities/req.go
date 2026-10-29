package entities

import (
	"encoding/json"
	"time"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"gorm.io/datatypes"
)

type DocumentResponse struct {
	DocumentID    string    `json:"documentID"`
	Title         string    `json:"title"`
	CollabVersion int       `json:"collabVersion"`
	OwnerID       string    `json:"ownerID"`
	Version       int       `json:"version"`
	Category      string    `json:"category"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy"`
}

type DocumentDetails struct {
	DocumentID    string         `json:"documentID"`
	Title         string         `json:"title"`
	Content       datatypes.JSON `json:"content"`
	CollabVersion int            `json:"collabVersion"`
	OwnerID       string         `json:"ownerID"`
	Version       int            `json:"version"`
	Category      string         `json:"category"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     *time.Time     `json:"deletedAt"`
	UpdatedBy     string         `json:"updatedBy"`
}

type DocumentReq struct {
	Title         string          `json:"title"`
	Content       json.RawMessage `json:"content"`
	Category      string          `json:"category"`
	CollabVersion *int            `json:"collabVersion"`
}

type DocumentSuccessRes struct {
	DocumentID string    `json:"documentID"`
	Title      string    `json:"title"`
	Version    int       `json:"version"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UpdatedBy  string    `json:"updatedBy"`
}
type RestoreVersionReq struct {
	ExpectedCollabVersion *int `json:"expectedCollabVersion"`
}

type AddUserRoleReq struct {
	Email string      `json:"email"`
	Role  enums.Roles `json:"role"`
}

type VersionsList struct {
	ID         string `json:"id"`
	DocumentID string `json:"docID"`
	Version    int    `json:"version"`
	Title      string `json:"title"`
	// Content            string    `json:"content"`
	ModifiedByUserID   string    `json:"modifiedByUserID"`
	ModifiedByUserName *string   `json:"modifiedByUserName"`
	CreatedAt          time.Time `json:"createdAt"`
}

type MemberDetailRes struct {
	UserID     string
	DocumentID string
	Role       enums.Roles
}

type UpdateRoleReq struct {
	Role enums.Roles `json:"role"`
}

type WSMessage struct {
	Type  string   `json:"type"`
	Data  EditData `json:"data,omitempty"`
	Error *WSError `json:"error,omitempty"`
}

type EditData struct {
	BaseCollabVersion int             `json:"baseCollabVersion"`
	Content           json.RawMessage `json:"content"`
}

type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BroadcastMessage struct {
	Type string        `json:"type"`
	Data BroadcastData `json:"data,omitempty"`
}

type BroadcastData struct {
	DocumentID    string         `json:"document_id"`
	CollabVersion int            `json:"collabVersion"`
	Content       datatypes.JSON `json:"content"`
	Title         string         `json:"title"`
}
