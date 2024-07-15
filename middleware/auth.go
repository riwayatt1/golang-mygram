package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/riwayatt1/golang-mygram/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is missing")
			c.Abort()
			return
		}

		// Expecting "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			utils.ErrorResponse(c, http.StatusInternalServerError, "JWT secret key is not set")
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			// Provide detailed error messages for debugging
			var errMsg string
			if err == jwt.ErrSignatureInvalid {
				errMsg = "Invalid token signature"
			} else {
				errMsg = "Invalid token"
			}
			utils.ErrorResponse(c, http.StatusUnauthorized, errMsg)
			c.Abort()
			return
		}

		if !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token is not valid")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Set("userID", uint(userID))
		c.Next()
	}
}

func AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		resourceIDStr := c.Param("userId")
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		if userID.(uint) != uint(resourceID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
