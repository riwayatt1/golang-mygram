package models

import (
	"time"
)

// SocialMedia represents the social media model
type SocialMedia struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	SocialMediaURL string    `gorm:"not null" json:"social_media_url"`
	UserID         uint      `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SocialMediaInput represents the input data for creating a social media record
type SocialMediaInput struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaURL string `json:"social_media_url" binding:"required"`
}
