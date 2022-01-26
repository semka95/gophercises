package repository

import (
	"context"
	"database/sql"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}


func (s sqliteRepository) GetByID(ctx context.Context, id string) (*PhoneNumber, error) {
	panic("implement me")
}

func (s sqliteRepository) Update(ctx context.Context, pn *PhoneNumber) error {
	panic("implement me")
}

func (s sqliteRepository) Store(ctx context.Context, pn *PhoneNumber) error {
	panic("implement me")
}

func (s sqliteRepository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}