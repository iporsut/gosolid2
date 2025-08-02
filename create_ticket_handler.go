package main

import (
	"fmt"
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

type CreateTicketHandler struct {
	db *gorm.DB // Assuming you're using GORM for database operations
}

func (h *CreateTicketHandler) Handler(c *gin.Context) {
	h.db.Transaction(func(tx *gorm.DB) error {
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

		var event Event
		if err := h.db.First(&event, eventID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		if req.Quantity > event.RemainingTickets {
			c.JSON(400, gin.H{"error": "Not enough tickets available"})
			return
		}

		event.RemainingTickets -= req.Quantity
		if err := tx.Save(&event).Error; err != nil {
			c.JSON(500, gin.H{"message": "Failed to update event tickets", "error": err.Error()})
			return
		}

		ticket := Ticket{
			EventID:      event.ID,
			Quantity:     req.Quantity,
			CustomerName: req.CustomerName,
			BookedAt:     time.Now(),
		}

		if err := tx.Create(&ticket).Error; err != nil {
			c.JSON(500, gin.H{"message": "Failed to book tickets", "error": err.Error()})
			return
		}

		c.JSON(201, CreateTicketResponse{
			ID:           ticket.ID,
			EventID:      ticket.EventID,
			Quantity:     ticket.Quantity,
			BookedAt:     ticket.BookedAt,
			CustomerName: ticket.CustomerName,
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		})
	})
}
