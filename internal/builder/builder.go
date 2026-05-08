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
	domainInterfaces []domain.DomainIF
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
		domainInterfaces: make([]domain.DomainIF, 0),
	}, nil
}

func (b *Builder) AddDomainInterface(entity domain.DomainIF) *Builder {
	b.domainInterfaces = append(b.domainInterfaces, entity)
	return b
}

func (b *Builder) Build(ctx context.Context) error {
	b.logger.Info(ctx, "Building domains...")
	defer b.logger.Info(ctx, "Building domains finished")

	api := chi.NewRouter()

	

	for _, entity := range b.domainInterfaces {
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
		defer b.logger.Info(rCtx, "Health check request finished.")
		for _, domain := range b.domainInterfaces {
			if err := domain.Health(rCtx); err != nil {
				b.logger.Errorf(rCtx, "Health check failed for domain: %v", err)
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func (b *Builder) Close(ctx context.Context) error {
	b.logger.Info(ctx, "Closing domains...")
	defer b.logger.Info(ctx, "Closing domains finished")
	for _, domain := range b.domainInterfaces {
		if err := domain.Close(ctx); err != nil {
			b.logger.Warnf(ctx, "Error closing entity: %v", err)
		}
	}
	return nil
}
