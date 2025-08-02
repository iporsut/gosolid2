package main

import (
	"crypto/rsa"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type JWTTokenCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginHandler struct {
	privateKey *rsa.PrivateKey
}

func (h *LoginHandler) Handler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Here you would typically check the username and password against your database
	if req.Username == "admin" && req.Password == "password" {
		// Generate JWT token (this is a placeholder, implement your JWT generation logic)
		t := jwt.NewWithClaims(jwt.SigningMethodRS256, JWTTokenCustomClaims{
			UserID: "12345",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "myapp",
				Subject:   "user",
				Audience:  jwt.ClaimStrings{"myapp_users"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})

		token, err := t.SignedString(h.privateKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(200, LoginResponse{Token: token})
	}
}
