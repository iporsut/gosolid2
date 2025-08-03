package main

import (
	"gorm.io/gorm"
)

type eventRepository struct {
	tx *gorm.DB
}

func NewEventRepository(tx *gorm.DB) *eventRepository {
	return &eventRepository{tx: tx}
}

func (r *eventRepository) GeByID(id uint) (*Event, error) {
	var event Event
	if err := r.tx.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Save(event *Event) error {
	if err := r.tx.Save(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) Create(event *Event) error {
	if err := r.tx.Create(event).Error; err != nil {
		return err
	}
	return nil
}
