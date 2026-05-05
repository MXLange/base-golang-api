package builder

import (
	"context"
	"net/http"

	"github.com/MXLange/go-model/internal/domain"
	"github.com/MXLange/go-model/internal/errors"
	"github.com/MXLange/go-model/internal/logger"
	"github.com/go-chi/chi/v5"
)

type Builder struct {
	logger         logger.LoggerIF
	r              *chi.Mux
	domainEntities []domain.EntityIF
}

func New(r *chi.Mux, logger logger.LoggerIF) (*Builder, error) {
	if r == nil {
		return nil, errors.ErrNilMux
	}

	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	return &Builder{
		r:              r,
		logger:         logger,
		domainEntities: make([]domain.EntityIF, 0),
	}, nil
}

func (b *Builder) AddDomainEntity(entity domain.EntityIF) *Builder {
	b.domainEntities = append(b.domainEntities, entity)
	return b
}

func (b *Builder) Build(ctx context.Context) error {
	b.logger.Info(ctx, "Building domain entities...")
	defer b.logger.Info(ctx, "Finished building domain entities.")

	api := chi.NewRouter()

	

	for _, entity := range b.domainEntities {
		if err := entity.Build(ctx, api); err != nil {
			return err
		}
	}

	prefix := "/api/v1"

	b.r.Mount(prefix, api)

	api.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		appError := errors.New(http.StatusMethodNotAllowed).WithError(errors.NewError("route", "method not allowed"))
		b.logger.Warnf(r.Context(), "Method not allowed: %s %s", r.Method, r.URL.Path)
		appError.WriteResponse(w)
	})
	
	api.NotFound(func(w http.ResponseWriter, r *http.Request) {
		appError := errors.New(http.StatusNotFound).WithError(errors.NewError("route", "not found"))
		b.logger.Warnf(r.Context(), "Not found: %s %s", r.Method, r.URL.Path)
		appError.WriteResponse(w)
	})


	b.setRoutes()
	return nil
}

func (b *Builder) setRoutes() {
	b.r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		b.logger.Info(rCtx, "Received health check request.")
		defer b.logger.Info(rCtx, "Finished processing health check request.")
		for _, entity := range b.domainEntities {
			if err := entity.Health(rCtx); err != nil {
				b.logger.Errorf(rCtx, "Health check failed for entity: %v", err)
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func (b *Builder) Close(ctx context.Context) error {
	b.logger.Info(ctx, "Closing domain entities...")
	defer b.logger.Info(ctx, "Finished closing domain entities.")
	for _, entity := range b.domainEntities {
		if err := entity.Close(ctx); err != nil {
			b.logger.Warnf(ctx, "Error closing entity: %v", err)
		}
	}
	return nil
}
