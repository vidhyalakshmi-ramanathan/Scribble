package entities

import (
	"encoding/json"
	"time"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
)

type DocumentReq struct {
	Title    string          `json:"title"`
	Content  json.RawMessage `json:"content"`
	Category string          `json:"category"`
}

type AddUserRoleReq struct {
	Email string      `json:"email"`
	Role  enums.Roles `json:"role"`
}

type VersionsList struct {
	ID         string `json:"id"`
	DocumentID string `json:"docID"`
	VersionID  int    `json:"versionID"`
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
	BaseVersion int    `json:"base_version"`
	Operation   string `json:"operation"`
	Position    int    `json:"position"`
	Text        string `json:"text"`
	Length      int    `json:"length"`
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
	DocumentID string `json:"document_id"`
	Version    int    `json:"version"`
	Content    string `json:"content"`
	Title      string `json:"title"`
}
