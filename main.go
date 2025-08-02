package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JWTTokenCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Event struct {
	gorm.Model
	Name             string    `gorm:"not null"`
	Description      string    `gorm:"not null"`
	NumberOfTickets  int       `gorm:"not null"`
	RemainingTickets int       `gorm:"not null"` // Remaining tickets after booking
	StartDateTime    time.Time `gorm:"not null"`
	Duration         int       `gorm:"not null"` // Duration in minutes
}

type Ticket struct {
	gorm.Model
	EventID      uint      `gorm:"not null"`
	Event        Event     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	Quantity     int       `gorm:"not null"` // Number of tickets booked
	BookedAt     time.Time `gorm:"not null"` // Timestamp when the ticket was booked
	CustomerName string    `gorm:"not null"` // Name of the customer who booked the ticket
}

type CreateEventRequest struct {
	Name            string    `json:"name" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	NumberOfTickets int       `json:"number_of_tickets" binding:"required"`
	StartDateTime   time.Time `json:"start_date_time" binding:"required"` // ISO 8601 format
	Duration        int       `json:"duration" binding:"required"`        // Duration in minutes
}

type CreateEventResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	NumberOfTickets  int       `json:"number_of_tickets"`
	StartDateTime    time.Time `json:"start_date_time"`
	Duration         int       `json:"duration"` // Duration in minutes
	RemainingTickets int       `json:"remaining_tickets"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func main() {
	// Gorm database connection to be established here
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=eventdb port=5432 sslmode=disable"))
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(
		&Event{},
		&Ticket{},
	); err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	privateKeyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("failed to read private key file: %v", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	publicKeyBytes, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalf("failed to read public key file: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	authorizeMiddleware := func(c *gin.Context) {
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
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		fmt.Println(token.Claims.(*JWTTokenCustomClaims)) // For debugging purposes

		c.Next()
	}

	e := gin.Default()

	// Create new event handler
	e.POST("/events", authorizeMiddleware, func(c *gin.Context) {
		var req CreateEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "invalid request", "error": err.Error()})
			return
		}

		event := Event{
			Name:             req.Name,
			Description:      req.Description,
			NumberOfTickets:  req.NumberOfTickets,
			StartDateTime:    req.StartDateTime,
			Duration:         req.Duration,
			RemainingTickets: req.NumberOfTickets, // Initially all tickets are available
		}

		if err := db.Create(&event).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create event"})
			return
		}
		c.JSON(201, CreateEventResponse{
			ID:               event.ID,
			Name:             event.Name,
			Description:      event.Description,
			NumberOfTickets:  event.NumberOfTickets,
			StartDateTime:    event.StartDateTime,
			Duration:         event.Duration,
			RemainingTickets: event.RemainingTickets,
			CreatedAt:        event.CreatedAt,
			UpdatedAt:        event.UpdatedAt,
		})
	})

	e.POST("/events/:id/tickets", authorizeMiddleware, func(c *gin.Context) {
		eventID := c.Param("id")
		if eventID == "" {
			c.JSON(400, gin.H{"message": "event ID is required"})
			return
		}
		var eventIDUint uint
		if _, err := fmt.Sscanf(eventID, "%d", &eventIDUint); err != nil {
			c.JSON(400, gin.H{"message": "invalid event ID format"})
			return
		}

		var req struct {
			CustomerName string `json:"customer_name" binding:"required"`
			Quantity     int    `json:"quantity" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "invalid request", "error": err.Error()})
			return
		}

		var event Event
		if err := db.First(&event, eventID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		if req.Quantity > event.RemainingTickets {
			c.JSON(400, gin.H{"error": "Not enough tickets available"})
			return
		}

		event.RemainingTickets -= req.Quantity
		if err := db.Save(&event).Error; err != nil {
			c.JSON(500, gin.H{"message": "Failed to update event tickets", "error": err.Error()})
			return
		}

		ticket := Ticket{
			EventID:      event.ID,
			Quantity:     req.Quantity,
			CustomerName: req.CustomerName,
			BookedAt:     time.Now(),
		}

		if err := db.Create(&ticket).Error; err != nil {
			c.JSON(500, gin.H{"message": "Failed to book tickets", "error": err.Error()})
			return
		}

		type TicketResponse struct {
			ID           uint      `json:"id"`
			EventID      uint      `json:"event_id"`
			Quantity     int       `json:"quantity"`
			BookedAt     time.Time `json:"booked_at"`
			CustomerName string    `json:"customer_name"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
		}

		c.JSON(201, TicketResponse{
			ID:           ticket.ID,
			EventID:      ticket.EventID,
			Quantity:     ticket.Quantity,
			BookedAt:     ticket.BookedAt,
			CustomerName: ticket.CustomerName,
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		})
	})

	// Authentication handler generate token if user is valid
	e.POST("/auth/login", func(c *gin.Context) {
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

			token, err := t.SignedString(privateKey)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(200, gin.H{"token": token})
		}
	})

	if err := e.Run(":8080"); err != nil {
		log.Fatalf("failed to run the server: %v", err)
	}
}
