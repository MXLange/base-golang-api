package user

import (
	"context"
	"net/http"

	"github.com/MXLange/go-model/internal/errors"
)

// handlersIF defines the interface for the handler layer of the userlication.
type handlersIF interface {
	Create(w http.ResponseWriter, r *http.Request)
}

// servicesIF defines the interface for the service layer of the userlication.
type servicesIF interface {
	Health(ctx context.Context) error
	Create(ctx context.Context, name string) (int, *errors.AppError)
}

// repositoryIF defines the interface for the repository layer of the userlication. It includes a Create method that takes a name as input and returns an integer (presumably an ID) and a pointer to an UserError in case of an error.
type repositoryIF interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, name string) (int, *errors.AppError)
}
