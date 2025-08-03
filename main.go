package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	authorizeMiddleware := AuthorizeMiddlewareHandler{publicKey: publicKey}

	e := gin.Default()

	// Create new event handler
	e.POST("/events", authorizeMiddleware.Handler, func(c *gin.Context) {
		h := CreateEventHandler{
			db:           db,
			eventService: NewEventService(NewEventRepository()),
		}
		h.Handler(c)
	})

	e.POST("/events/:id/tickets", authorizeMiddleware.Handler, func(c *gin.Context) {
		h := CreateTicketHandler{
			db: db,
			ticketService: NewTicketService(
				NewTicketRepository(),
				NewEventRepository(),
			),
		}
		h.Handler(c)
	})

	// Authentication handler generate token if user is valid
	e.POST("/auth/login", func(c *gin.Context) {
		h := LoginHandler{privateKey: privateKey}
		h.Handler(c)
	})

	if err := e.Run(":8080"); err != nil {
		log.Fatalf("failed to run the server: %v", err)
	}
}
