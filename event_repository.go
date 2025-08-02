package main

import (
	"gorm.io/gorm"
)

type EventRepository struct {
	tx *gorm.DB
}

func NewEventRepository(tx *gorm.DB) *EventRepository {
	return &EventRepository{tx: tx}
}

func (r *EventRepository) GeByID(id uint) (*Event, error) {
	var event Event
	if err := r.tx.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Save(event *Event) error {
	if err := r.tx.Save(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *EventRepository) Create(event *Event) error {
	if err := r.tx.Create(event).Error; err != nil {
		return err
	}
	return nil
}
