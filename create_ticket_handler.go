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

type CreateTicketHandler struct {
	db *gorm.DB // Assuming you're using GORM for database operations
}

type TicketService struct {
	tx *gorm.DB
}

func NewTicketService(tx *gorm.DB) *TicketService {
	return &TicketService{tx: tx}
}

type Error struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Message: %s, StatusCode: %d, Err: %v", e.Message, e.StatusCode, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (s *TicketService) CreateTicket(eventID uint, req CreateTicketRequest) (*CreateTicketResponse, error) {
	var event Event
	if err := s.tx.First(&event, eventID).Error; err != nil {
		return nil, &Error{Message: "event not found", StatusCode: http.StatusNotFound, Err: err}
	}

	if req.Quantity > event.RemainingTickets {
		return nil, &Error{Message: "not enough tickets available", StatusCode: http.StatusBadRequest}
	}

	event.RemainingTickets -= req.Quantity
	if err := s.tx.Save(&event).Error; err != nil {
		return nil, &Error{Message: "failed to update event tickets", StatusCode: http.StatusInternalServerError, Err: err}
	}

	ticket := Ticket{
		EventID:      event.ID,
		Quantity:     req.Quantity,
		CustomerName: req.CustomerName,
		BookedAt:     time.Now(),
	}

	if err := s.tx.Create(&ticket).Error; err != nil {
		return nil, &Error{Message: "failed to book tickets", StatusCode: http.StatusInternalServerError, Err: err}
	}

	return &CreateTicketResponse{
		ID:           ticket.ID,
		EventID:      ticket.EventID,
		Quantity:     ticket.Quantity,
		BookedAt:     ticket.BookedAt,
		CustomerName: ticket.CustomerName,
		CreatedAt:    ticket.CreatedAt,
		UpdatedAt:    ticket.UpdatedAt,
	}, nil
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
		s := NewTicketService(tx)
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
