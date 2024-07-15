package models

import (
	"time"
)

// Comment represents the comment model
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `gorm:"not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentInput represents the input data for creating a comment
type CommentInput struct {
	Message string `json:"message" binding:"required"`
}
