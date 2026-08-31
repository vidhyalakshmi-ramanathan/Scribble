package entities

import "gorm.io/gorm"

type Comments struct {
	gorm.Model
	DocumentID string
	CommentBy  string
	Position   string
	CommentTo  string
}
