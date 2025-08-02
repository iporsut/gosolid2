package main

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

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
	eventRepo := NewEventRepository(s.tx)

	event, err := eventRepo.GeByID(eventID)
	if err != nil {
		return nil, &Error{Message: "event not found", StatusCode: http.StatusNotFound, Err: err}
	}

	if req.Quantity > event.RemainingTickets {
		return nil, &Error{Message: "not enough tickets available", StatusCode: http.StatusBadRequest}
	}

	event.RemainingTickets -= req.Quantity
	if err := eventRepo.Save(event); err != nil {
		return nil, &Error{Message: "failed to update event tickets", StatusCode: http.StatusInternalServerError, Err: err}
	}

	ticketRepo := NewTicketRepository(s.tx)

	ticket := Ticket{
		EventID:      event.ID,
		Quantity:     req.Quantity,
		CustomerName: req.CustomerName,
		BookedAt:     time.Now(),
	}

	if err := ticketRepo.Create(&ticket); err != nil {
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
