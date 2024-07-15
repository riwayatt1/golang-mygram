package models

import (
	"time"
)

// Photo represents the photo model
type Photo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" validate:"required"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url" validate:"required,url"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PhotoInput represents the input data for creating a photo
type PhotoInput struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" binding:"required"`
}
