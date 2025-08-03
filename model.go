package main

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name             string    `gorm:"not null"`
	Description      string    `gorm:"not null"`
	NumberOfTickets  int       `gorm:"not null"`
	RemainingTickets int       `gorm:"not null"` // Remaining tickets after booking
	StartDateTime    time.Time `gorm:"not null"`
	Duration         int       `gorm:"not null"` // Duration in minutes
}

var (
	ErrNotEnoughTickets = errors.New("not enough tickets available")
)

func (e *Event) NewOrderTicket(quantity int, customerName string) (*Ticket, error) {
	if quantity <= 0 || quantity > e.RemainingTickets {
		return nil, errors.New("invalid ticket quantity")
	}

	ticket := &Ticket{
		EventID:      e.ID,
		Quantity:     quantity,
		BookedAt:     time.Now(),
		CustomerName: customerName,
	}

	if e.RemainingTickets < quantity {
		return nil, ErrNotEnoughTickets
	}

	e.RemainingTickets -= quantity

	return ticket, nil
}

type Ticket struct {
	gorm.Model
	EventID      uint      `gorm:"not null"`
	Event        Event     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	Quantity     int       `gorm:"not null"` // Number of tickets booked
	BookedAt     time.Time `gorm:"not null"` // Timestamp when the ticket was booked
	CustomerName string    `gorm:"not null"` // Name of the customer who booked the ticket
}
