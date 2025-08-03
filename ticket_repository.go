package main

import (
	"context"
	"errors"
)

type ticketRepository struct {
}

func NewTicketRepository() *ticketRepository {
	return &ticketRepository{}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *Ticket) error {
	tx := TxFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}

	if err := tx.WithContext(ctx).Create(ticket).Error; err != nil {
		return err
	}
	return nil
}
