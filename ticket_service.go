package main

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type newTicketRepositoryFunc func(tx *gorm.DB) TicketRepository
type newEventRepositoryFunc func(tx *gorm.DB) EventRepository

type ticketService struct {
	tx                  *gorm.DB
	newTicketRepository newTicketRepositoryFunc
	newEventRepository  newEventRepositoryFunc
}

func NewTicketService(tx *gorm.DB, newTicketRepository newTicketRepositoryFunc, newEventRepository newEventRepositoryFunc) *ticketService {
	return &ticketService{
		tx:                  tx,
		newTicketRepository: newTicketRepository,
		newEventRepository:  newEventRepository,
	}
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

type TicketRepository interface {
	Create(ticket *Ticket) error
}

type EventRepository interface {
	GeByID(id uint) (*Event, error)
	Save(event *Event) error
	Create(event *Event) error
}

func (s *ticketService) CreateTicket(eventID uint, req CreateTicketRequest) (*CreateTicketResponse, error) {
	eventRepo := s.newEventRepository(s.tx)

	event, err := eventRepo.GeByID(eventID)
	if err != nil {
		return nil, &Error{Message: "event not found", StatusCode: http.StatusNotFound, Err: err}
	}

	ticket, err := event.NewOrderTicket(req.Quantity, req.CustomerName)
	if err != nil {
		if err == ErrNotEnoughTickets {
			return nil, &Error{Message: "not enough tickets available", StatusCode: http.StatusBadRequest, Err: err}
		}
		return nil, &Error{Message: "failed to create ticket", StatusCode: http.StatusBadRequest, Err: err}
	}

	if err := eventRepo.Save(event); err != nil {
		return nil, &Error{Message: "failed to update event tickets", StatusCode: http.StatusInternalServerError, Err: err}
	}

	ticketRepo := s.newTicketRepository(s.tx)

	if err := ticketRepo.Create(ticket); err != nil {
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
