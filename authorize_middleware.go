package main

import (
	"crypto/rsa"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthorizeMiddlewareHandler struct {
	publicKey *rsa.PublicKey
}

func (h *AuthorizeMiddlewareHandler) Handler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		c.JSON(401, gin.H{"error": "Authorization header must start with Bearer"})
		c.Abort()
		return
	}

	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return h.publicKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	fmt.Println(token.Claims.(*JWTTokenCustomClaims)) // For debugging purposes

	c.Next()
}
