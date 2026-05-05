package app

import (
	"io"
	"net/http"

	"github.com/MXLange/go-model/internal/domain/app/dto"
	"github.com/MXLange/go-model/internal/errors"
	jsonschema "github.com/MXLange/go-model/internal/json_schema"
	"github.com/MXLange/go-model/internal/logger"
)

type handlers struct {
	name    string
	logger  logger.LoggerIF
	service servicesIF
}

func NewHandlers(name string, services servicesIF, logger logger.LoggerIF) (handlersIF, error) {
	if services == nil {
		return nil, errors.ErrNilService
	}

	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	return &handlers{
		name:    name,
		logger:  logger,
		service: services,
	}, nil
}

func (h *handlers) GetName() string {
	return h.name
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof(r.Context(), "%s handler - received request to create.", h.name)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf(r.Context(), "%s handler - failed to read request body: %v", h.name, err)
		appError := errors.New(http.StatusBadRequest).WithError(errors.NewError("", "invalid request body"))
		appError.WriteResponse(w)
		return
	}

	schemaValidator := jsonschema.New(dto.CreateInSchema)
	if err := schemaValidator.Validate(string(body)); err != nil {
		h.logger.Errorf(r.Context(), "%s handler - request body failed schema validation: %v", h.name, err)
		err.WriteResponse(w)
		return
	}

	w.Write([]byte("hello from " + h.name + " handler"))
}
