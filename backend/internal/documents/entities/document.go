package entities

import (
	"time"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Documents defines the data being store in the DB about a particular document and its ownwer
type Documents struct {
	DocumentID string         `gorm:"primaryKey" json:"docID"`
	Title      string         `gorm:"not null" json:"title"`
	Content    datatypes.JSON `gorm:"type:jsonb;not null" json:"content"` // ProseMirror doc JSON
	SearchText string         `gorm:"type:text" json:"-"`                 // derived, plain text
	SchemaVer  int            `gorm:"not null;default:1" json:"-"`        // content schema version
	OwnerID    string         `gorm:"not null" json:"ownerID"`
	Version    int            `gorm:"not null;default:0" json:"version"`
	Category   string         `gorm:"not null;default:General" json:"category"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	UpdatedBy  string         `gorm:"not null" json:"updatedBY"`
}

// Permissions defines Who can access the document, and what can they do
type Permissions struct {
	ID         string      `gorm:"primaryKey"`
	UserID     string      `gorm:"not null;uniqueIndex:idx_document_user"`
	DocumentID string      `gorm:"not null;uniqueIndex:idx_document_user"`
	Role       enums.Roles `gorm:"type:varchar(20);not null"`
}

type Version struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	DocumentID string         `gorm:"not null;uniqueIndex:idx_document_version" json:"docID"`
	VersionID  int            `gorm:"not null;uniqueIndex:idx_document_version" json:"versionID"`
	Title      string         `gorm:"not null" json:"title"`
	Content    datatypes.JSON `gorm:"type:jsonb;not null" json:"content"`
	ModifiedBy string         `gorm:"not null" json:"modifiedBy"`
	CreatedAt  time.Time      `json:"createdAt"`
}
