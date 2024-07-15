package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riwayatt1/golang-mygram/config"
	"github.com/riwayatt1/golang-mygram/models"
	"github.com/riwayatt1/golang-mygram/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Username already exists")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the email already exists
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email already exists")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create new user
	newUser := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Age:       input.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "created", newUser)
}

func LoginUser(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, "ok", gin.H{"token": token})
}

func GetUserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "ok", user)
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	var updateUser models.User
	if err := config.DB.First(&updateUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reflect on the input struct to dynamically update fields
	inputVal := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	for i := 0; i < inputVal.NumField(); i++ {
		field := inputVal.Field(i)
		fieldName := inputType.Field(i).Name

		// Check if the field has a zero value (meaning it wasn't set in the input)
		if !field.IsZero() {
			updateField := reflect.ValueOf(&updateUser).Elem().FieldByName(fieldName)
			if updateField.IsValid() && updateField.CanSet() {
				updateField.Set(field)
			}
		}
	}

	// Validate the updated user struct
	if err := utils.ValidateStruct(updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated user
	if err := config.DB.Save(&updateUser).Error; err != nil {
		// Check if the error is a unique constraint violation
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "idx_users_username") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			} else if strings.Contains(err.Error(), "idx_users_email") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updateUser})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	var deleteUser models.User
	if err := config.DB.First(&deleteUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&deleteUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your account has been sucessfully deleted"})
}
