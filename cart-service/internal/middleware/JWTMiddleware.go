package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWTMiddleware is a middleware function to authenticate JWT tokens
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		// Extract the token from "Bearer <token>"
		tokenString := extractToken(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// You should store your secret key in a secure location
			// Here, we are using a hardcoded key for demonstration purposes
			return []byte("your_secret_key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		// Extract data from the token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract data from claims
			fmt.Println(claims)
			if userUuid, ok := claims["uuid"].(string); ok {
				// Set extracted data to the context for further use
				c.Set("user_uuid", userUuid)
				c.Set("token", tokenString)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_uuid type in token claims"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Token is valid, proceed with the request
		c.Next()
	}
}

// extractToken extracts the token from "Bearer <token>"
func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
