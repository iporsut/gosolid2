package main

import (
	"gorm.io/gorm"
)

type ticketRepository struct {
	tx *gorm.DB
}

func NewTicketRepository(tx *gorm.DB) *ticketRepository {
	return &ticketRepository{tx: tx}
}

func (r *ticketRepository) Create(ticket *Ticket) error {
	if err := r.tx.Create(ticket).Error; err != nil {
		return err
	}
	return nil
}
