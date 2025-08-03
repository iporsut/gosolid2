package main

import (
	"net/http"

	"gorm.io/gorm"
)

type eventService struct {
	tx                 *gorm.DB
	newEventRepository newEventRepositoryFunc
}

func NewEventService(tx *gorm.DB, newEventRepository newEventRepositoryFunc) *eventService {
	return &eventService{tx: tx, newEventRepository: newEventRepository}
}

func (s *eventService) CreateEvent(event *Event) error {
	eventRepo := s.newEventRepository(s.tx)
	if err := eventRepo.Create(event); err != nil {
		return &Error{Message: "failed to create event", StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}
