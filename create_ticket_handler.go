package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTicketRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required"`
}

type CreateTicketResponse struct {
	ID           uint      `json:"id"`
	EventID      uint      `json:"event_id"`
	Quantity     int       `json:"quantity"`
	BookedAt     time.Time `json:"booked_at"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type newTicketServiceFunc func(tx *gorm.DB) TicketService

type CreateTicketHandler struct {
	db               *gorm.DB
	newTicketService newTicketServiceFunc
}

func (h *CreateTicketHandler) Handler(c *gin.Context) {
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
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		s := h.newTicketService(tx)
		response, err := s.CreateTicket(eventIDUint, req)
		if err != nil {
			return err
		}

		c.JSON(201, response)
		return nil
	})

	if err != nil {
		c.Error(err)
		if customErr, ok := err.(*Error); ok {
			c.JSON(customErr.StatusCode, gin.H{"message": customErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "error": err.Error()})
		}
		return
	}
}
