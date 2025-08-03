package main

import (
	"context"
	"fmt"
	"net/http"
)

type ticketService struct {
	ticketRepository TicketRepository
	eventRepository  EventRepository
}

func NewTicketService(ticketRepository TicketRepository, eventRepository EventRepository) *ticketService {
	return &ticketService{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
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

func (s *ticketService) CreateTicket(ctx context.Context, eventID uint, req CreateTicketRequest) (*CreateTicketResponse, error) {
	event, err := s.eventRepository.GeByID(ctx, eventID)
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

	if err := s.eventRepository.Save(ctx, event); err != nil {
		return nil, &Error{Message: "failed to update event tickets", StatusCode: http.StatusInternalServerError, Err: err}
	}

	if err := s.ticketRepository.Create(ctx, ticket); err != nil {
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
