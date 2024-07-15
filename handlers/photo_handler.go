package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riwayatt1/golang-mygram/config"
	"github.com/riwayatt1/golang-mygram/models"
	"github.com/riwayatt1/golang-mygram/utils"
)

// CreatePhoto handles the creation of a new photo
func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	photo.UserID = userID

	config.DB.Create(&photo)
	utils.SuccessResponse(c, "created", photo)
}

func GetAllPhotos(c *gin.Context) {
	// Get userID from context (set by your authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var photos []models.Photo
	if err := config.DB.Preload("User").Where("user_id = ?", userID).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transforming the response to include user information
	var response []gin.H
	for _, photo := range photos {
		userData := gin.H{
			"email":    photo.User.Email,
			"username": photo.User.Username,
		}
		photoData := gin.H{
			"id":         photo.ID,
			"title":      photo.Title,
			"caption":    photo.Caption,
			"photo_url":  photo.PhotoURL,
			"user_id":    photo.UserID,
			"created_at": photo.CreatedAt,
			"updated_at": photo.UpdatedAt,
			"user":       userData,
		}
		response = append(response, photoData)
	}

	utils.SuccessResponse(c, "ok", response)
}

// GetPhoto retrieves a photo by its ID
func GetPhoto(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo
	if err := config.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	utils.SuccessResponse(c, "ok", photo)
}

// UpdatePhoto updates a photo by its ID
func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo
	if err := config.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&photo)
	utils.SuccessResponse(c, "ok", photo)
}

// DeletePhoto deletes a photo by its ID
func DeletePhoto(c *gin.Context) {
	// Get photo ID from URL parameter
	id := c.Param("id")

	// Get userID from context (set by your authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// Fetch the photo from the database
	var photo models.Photo
	if err := config.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Check if the authenticated user owns the photo
	if photo.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this photo"})
		return
	}

	// Delete the photo
	if err := config.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, "ok", gin.H{"message": "Your photo has been successfully deleted"})
}
