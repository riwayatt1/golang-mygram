package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riwayatt1/golang-mygram/config"
	"github.com/riwayatt1/golang-mygram/models"
)

// CreateComment handles the creation of a new comment
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	comment.UserID = userID

	config.DB.Create(&comment)
	c.JSON(http.StatusCreated, gin.H{"data": comment})
}

// GetComment retrieves a comment by its ID
func GetComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := config.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// UpdateComment updates a comment by its ID
func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := config.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&comment)
	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// DeleteComment deletes a comment by its ID
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Comment{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
