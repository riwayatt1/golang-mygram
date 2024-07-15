package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, status string, data interface{}) {
	statusCode := http.StatusOK // Default to 200 OK
	switch status {
	case "ok":
		statusCode = http.StatusOK
	case "created":
		statusCode = http.StatusCreated
	}

	// if IsGinH(data) {
	// 	responseData := data
	// 	c.JSON(statusCode, Response{
	// 		Status: "success",
	// 		Data:   responseData,
	// 	})
	// 	return
	// }
	// responseData := RemoveFieldsFromMap(data, "CreatedAt", "UpdatedAt", "Password")

	c.JSON(statusCode, Response{
		Status: "success",
		Data:   data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, Response{
		Status: "error",
		Data: gin.H{
			"error": errorMessage,
		},
	})
}
