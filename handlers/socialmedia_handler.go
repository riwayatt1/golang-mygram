package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riwayatt1/golang-mygram/config"
	"github.com/riwayatt1/golang-mygram/models"
)

// CreateSocialMedia handles the creation of a new social media record
func CreateSocialMedia(c *gin.Context) {
	var socialMedia models.SocialMedia
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	socialMedia.UserID = userID

	config.DB.Create(&socialMedia)
	c.JSON(http.StatusCreated, gin.H{"data": socialMedia})
}

// GetSocialMedia retrieves a social media record by its ID
func GetSocialMedia(c *gin.Context) {
	id := c.Param("id")
	var socialMedia models.SocialMedia
	if err := config.DB.Where("id = ?", id).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": socialMedia})
}

// UpdateSocialMedia updates a social media record by its ID
func UpdateSocialMedia(c *gin.Context) {
	id := c.Param("id")
	var socialMedia models.SocialMedia
	if err := config.DB.Where("id = ?", id).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media record not found"})
		return
	}

	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&socialMedia)
	c.JSON(http.StatusOK, gin.H{"data": socialMedia})
}

// DeleteSocialMedia deletes a social media record by its ID
func DeleteSocialMedia(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.SocialMedia{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media record deleted successfully"})
}
