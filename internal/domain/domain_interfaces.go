package domain

import (
	"context"

	"github.com/go-chi/chi/v5"
)

type DomainIF interface {
	Build(ctx context.Context, r *chi.Mux) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error
}
