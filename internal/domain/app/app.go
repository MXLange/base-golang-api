package app

import (
	"context"
	"database/sql"

	"github.com/MXLange/go-model/internal/errors"
	"github.com/MXLange/go-model/internal/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	name       string
	logger     logger.LoggerIF
	handlers   handlersIF
	services   servicesIF
	repository repositoryIF
	db         *sql.DB
}

func NewApp(ctx context.Context, db *sql.DB, logger logger.LoggerIF) (*App, error) {

	if db == nil {
		return nil, errors.ErrNilDB
	}

	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	name := "App"

	logger.Infof(ctx, "Initializing %s application.", name)
	defer logger.Infof(ctx, "Finished initializing %s application.", name)

	repository, err := NewRepository(name, db, logger)
	if err != nil {
		return nil, err
	}

	services, err := NewServices(name, repository, logger)
	if err != nil {
		return nil, err
	}

	handlers, err := NewHandlers(name, services, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		name:       name,
		logger:     logger,
		handlers:   handlers,
		services:   services,
		repository: repository,
		db:         db,
	}, nil
}

func (a *App) GetServices() servicesIF {
	return a.services
}

func (a *App) GetName() string {
	return a.name
}

func (a *App) Build(ctx context.Context, r *chi.Mux) error {
	a.logger.Infof(ctx, "Building %s application routes.", a.name)
	defer a.logger.Infof(ctx, "Finished building %s application routes.", a.name)
	return newAppRoutes(ctx, r, a.handlers)
}

func (a *App) Health(ctx context.Context) error {
	a.logger.Infof(ctx, "%s application - received request to ping the database.", a.name)
	return a.services.Health(ctx)
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Infof(ctx, "Closing %s application.", a.name)
	defer a.logger.Infof(ctx, "Finished closing %s application.", a.name)
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return err
		}
	}
	return nil
}
