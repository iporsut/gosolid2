package main

import (
	"net/http"

	"gorm.io/gorm"
)

type EventService struct {
	tx *gorm.DB
}

func NewEventService(tx *gorm.DB) *EventService {
	return &EventService{tx: tx}
}

func (s *EventService) CreateEvent(event *Event) error {
	eventRepo := NewEventRepository(s.tx)
	if err := eventRepo.Create(event); err != nil {
		return &Error{Message: "failed to create event", StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}
