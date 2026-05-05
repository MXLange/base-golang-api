package app

import (
	"context"

	"github.com/MXLange/go-model/internal/errors"
	"github.com/go-chi/chi/v5"
)

func newAppRoutes(ctx context.Context, r chi.Router, handlers handlersIF) error {
	if handlers == nil {
		return errors.ErrNilHandler
	}

	if r == nil {
		return errors.ErrNilMux
	}

	r.Post("/app", handlers.Create)

	return nil
}
