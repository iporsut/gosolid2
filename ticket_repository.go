package main

import (
	"gorm.io/gorm"
)

type TicketRepository struct {
	tx *gorm.DB
}

func NewTicketRepository(tx *gorm.DB) *TicketRepository {
	return &TicketRepository{tx: tx}
}

func (r *TicketRepository) Create(ticket *Ticket) error {
	if err := r.tx.Create(ticket).Error; err != nil {
		return err
	}
	return nil
}
