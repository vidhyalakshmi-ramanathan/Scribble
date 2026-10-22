package entities

import (
	"time"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
)

// User defines the data being store in the DB about a particular user
type User struct {
	UserID   string `gorm:"primaryKey" json:"userID"`
	UserName string `gorm:"size:100;not null"  json:"username"`
	Email    string `gorm:"unique;not null"  json:"email"`
	Password string `gorm:"not null"`
}

type DocumentInvites struct {
	ID           string      `gorm:"primaryKey"`
	DocumentID   string      `gorm:"not null"`
	InvitedEmail string      `gorm:"not null"`
	Role         enums.Roles `gorm:"not null"`
	// Token        string    `gorm:"uniqueIndex"`
	InvitedBy string        `gorm:"not null"`
	Status    enums.Invites `gorm:"not null"`
	ExpiresAt time.Time     `gorm:"not null"`
	CreatedAt time.Time
}
