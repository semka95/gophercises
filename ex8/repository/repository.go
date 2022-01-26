package repository

import (
	"context"
	"time"
)

type PhoneNumber struct {
	ID int
	Number string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*PhoneNumber, error)
	Update(ctx context.Context, pn *PhoneNumber) error
	Store(ctx context.Context, pn *PhoneNumber) error
	Delete(ctx context.Context, id string) error
}
