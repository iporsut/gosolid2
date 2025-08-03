package main

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	gormpostgres "gorm.io/driver/postgres"
)

type fakeTicketRepository struct {
	TicketRepository
}

func (f *fakeTicketRepository) Create(ctx context.Context, ticket *Ticket) error {
	// Simulate ticket creation logic
	ticket.ID = 1 // Assign a fake ID for testing
	return nil
}

type fakeEventRepository struct {
	EventRepository
}

func (f *fakeEventRepository) GeByID(ctx context.Context, id uint) (*Event, error) {
	// Simulate fetching an event by ID
	return &Event{
		Model: gorm.Model{
			ID: id,
		},
		Name:             "Test Event",
		Description:      "This is a test event",
		NumberOfTickets:  100,
		RemainingTickets: 100,
		StartDateTime:    time.Now(),
		Duration:         120,
	}, nil
}

func (f *fakeEventRepository) Save(ctx context.Context, event *Event) error {
	// Simulate saving an event
	return nil
}

func TestTicketService(t *testing.T) {
	t.Run("CreateTicket", func(t *testing.T) {
		s := NewTicketService(&fakeTicketRepository{}, &fakeEventRepository{})
		ctx := context.Background()
		req := CreateTicketRequest{
			Quantity:     2,
			CustomerName: "John Doe",
		}
		eventID := uint(1)
		_, err := s.CreateTicket(ctx, eventID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestCreateTickcetService(t *testing.T) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer postgresContainer.Terminate(context.Background())

	dbURI, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := gorm.Open(gormpostgres.Open(dbURI))

	err = db.AutoMigrate(&Event{}, &Ticket{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Create Event
	event := &Event{
		Name:             "Test Event",
		Description:      "This is a test event",
		NumberOfTickets:  100,
		RemainingTickets: 100,
		StartDateTime:    time.Now(),
		Duration:         120,
	}
	err = db.Create(event).Error
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	s := NewTicketService(NewTicketRepository(), NewEventRepository())

	db.Transaction(func(tx *gorm.DB) error {
		ctx := WithTx(ctx, tx)

		req := CreateTicketRequest{
			Quantity:     2,
			CustomerName: "John Doe",
		}

		resp, err := s.CreateTicket(ctx, event.ID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Quantity != req.Quantity {
			t.Fatalf("expected quantity %d, got %d", req.Quantity, resp.Quantity)
		}
		if resp.CustomerName != req.CustomerName {
			t.Fatalf("expected customer name %s, got %s", req.CustomerName, resp.CustomerName)
		}
		if resp.EventID != event.ID {
			t.Fatalf("expected event ID %d, got %d", event.ID, resp.EventID)
		}

		var ticket Ticket

		err = tx.First(&ticket, resp.ID).Error
		if err != nil {
			t.Fatalf("failed to retrieve ticket: %v", err)
		}
		t.Logf("Ticket created successfully: %+v", ticket)

		return nil
	})

}
