package main

import (
	"context"
	"errors"
)

type eventRepository struct {
}

func NewEventRepository() *eventRepository {
	return &eventRepository{}
}

func (r *eventRepository) GeByID(ctx context.Context, id uint) (*Event, error) {
	tx := TxFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}

	var event Event
	if err := tx.WithContext(ctx).First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Save(ctx context.Context, event *Event) error {
	tx := TxFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}

	if err := tx.WithContext(ctx).Save(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) Create(ctx context.Context, event *Event) error {
	tx := TxFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}

	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		return err
	}
	return nil
}
