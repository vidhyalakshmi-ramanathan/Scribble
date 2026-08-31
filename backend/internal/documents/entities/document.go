package entities

import (
	"time"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"gorm.io/gorm"
)

// Documents defines the data being store in the DB about a particular document and its ownwer
type Documents struct {
	DocumentID string `gorm:"primaryKey"`
	Title      string `gorm:"not null"`
	Content    string `gorm:"type:text"` //DocURL
	OwnerID    string `gorm:"not null"`
	Version    int    `gorm:"not null;default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UpdatedBy  string         `gorm:"not null"`
}

// Permissions defines Who can access the document, and what can they do
type Permissions struct {
	ID         string      `gorm:"primaryKey"`
	UserID     string      `gorm:"not null;uniqueIndex:idx_document_user"`
	DocumentID string      `gorm:"not null;uniqueIndex:idx_document_user"`
	Role       enums.Roles `gorm:"type:varchar(20);not null"`
}

type Version struct {
	ID         string `gorm:"primaryKey"`
	DocumentID string `gorm:"not null;uniqueIndex:idx_document_version"`
	VersionID  int    `gorm:"not null;uniqueIndex:idx_document_version"`
	Content    string `gorm:"type:text;not null"` //DocURL
	ModifiedBy string `gorm:"not null"`
	CreatedAt  time.Time
}
