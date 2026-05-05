package app

import (
	"context"

	"github.com/MXLange/go-model/internal/errors"
	"github.com/MXLange/go-model/internal/logger"
)

type services struct {
	name       string
	logger     logger.LoggerIF
	repository repositoryIF
}

func NewServices(name string, repository repositoryIF, logger logger.LoggerIF) (servicesIF, error) {
	if repository == nil {
		return nil, errors.ErrNilRepository
	}

	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	return &services{
		name:       name,
		repository: repository,
		logger:     logger,
	}, nil
}

func (s *services) Health(ctx context.Context) error {
	s.logger.Infof(ctx, "%s service - received request to ping the database.", s.name)
	return s.repository.Ping(ctx)
}

func (s *services) Create(ctx context.Context, name string) (int, *errors.AppError) {
	s.logger.Infof(ctx, "%s service - received request to create.", s.name)
	return s.repository.Create(ctx, name)
}
