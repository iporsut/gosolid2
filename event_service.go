package main

import (
	"context"
	"net/http"
)

type eventService struct {
	eventRepository EventRepository
}

func NewEventService(eventRepository EventRepository) *eventService {
	return &eventService{eventRepository: eventRepository}
}

func (s *eventService) CreateEvent(ctx context.Context, event *Event) error {
	if err := s.eventRepository.Create(ctx, event); err != nil {
		return &Error{Message: "failed to create event", StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}
