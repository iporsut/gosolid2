package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

type CreateEventHandler struct {
	db *gorm.DB // Assuming you're using GORM for database operations
}

func (h *CreateEventHandler) Handler(c *gin.Context) {
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

	if err := h.db.Create(&event).Error; err != nil {
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
}
