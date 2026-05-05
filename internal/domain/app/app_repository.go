package app

import (
	"context"
	"database/sql"

	"github.com/MXLange/go-model/internal/errors"
	"github.com/MXLange/go-model/internal/logger"
)

type repository struct {
	name   string
	logger logger.LoggerIF
	db     *sql.DB
}

func NewRepository(name string, db *sql.DB, logger logger.LoggerIF) (repositoryIF, error) {
	if db == nil {
		return nil, errors.ErrNilDB
	}

	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	return &repository{
		name:   name,
		logger: logger,
		db:     db,
	}, nil
}

func (s *repository) Ping(ctx context.Context) error {
	s.logger.Infof(ctx, "%s repository - received request to ping the database.", s.name)
	return s.db.PingContext(ctx)
}

func (s *repository) Create(ctx context.Context, name string) (int, *errors.AppError) {
	s.logger.Infof(ctx, "%s repository - received request to create.", s.name)
	return 0, nil
}
