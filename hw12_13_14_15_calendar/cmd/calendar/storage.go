package main

import (
	"context"
	"errors"

	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	storageTypeMemory = "memory"
	storageTypeSQL    = "sql"
)

var undefinedStorageType = errors.New("undefined storage type")

func NewStorage(ctx context.Context, config Config) (app.Storage, error) {
	var storage app.Storage

	switch config.App.StorageType {
	case storageTypeMemory:
		storage = memorystorage.New()
	case storageTypeSQL:
		db, err := sqlx.Open("pgx", config.PostgreSQL.BuildDSN())
		if err != nil {
			return nil, err
		}
		ss := sqlstorage.New(db)
		if err := ss.Connect(ctx); err != nil {
			return nil, err
		}
		storage = ss
	default:
		return nil, undefinedStorageType
	}

	return storage, nil
}
