package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrNilLogger     = errors.New("logger is nil")
	ErrNilHandler    = errors.New("handler is nil")
	ErrNilMux        = errors.New("mux is nil")
	ErrNilService    = errors.New("service is nil")
	ErrNilRepository = errors.New("repository is nil")
	ErrNilDB         = errors.New("db is nil")
)

// AppError represents a structured error response that can be returned by an API. It includes an HTTP status code and an optional list of errors, where each error can have a field and a message. This structure allows for consistent error handling and clear communication of issues to API consumers.
type AppError struct {
	code   int   `json:"-"`
	errors []FieldError `json:"errors,omitempty"`
}

// FieldError represents a single error detail that can be included in the AppError's errors field. It contains an optional field name and an optional message describing the error. This allows for granular error reporting, where each error can specify which field is problematic and what the issue is.
type FieldError struct {
	field   string `json:"field,omitempty" example:"name"`
	message string `json:"message,omitempty" example:"name is required"`
}

// New creates a new AppError with the given code. The code parameter is required and should be a valid HTTP status code. The errors field is optional and can be set using the WithErrors method if needed.
func New(code int) *AppError {
	return &AppError{
		code: code,
	}
}

// NewError creates a new error with the given field and message. Both parameters are optional and can be set to nil if not needed.
func NewError(field, message string) *FieldError {
	return &FieldError{
		field:   field,
		message: message,
	}
}

// GetCode returns the HTTP status code associated with the AppError. This method allows API handlers to easily retrieve the status code when constructing HTTP responses based on the error.
func (ae *AppError) GetCode() int {
	return ae.code
}

// WithError appends a single error to the errors field of the AppError. This method allows for chaining and returns the modified AppError instance, enabling a fluent interface for building error responses.
func (ae *AppError) WithError(errors *FieldError) *AppError {
	ae.errors = append(ae.errors, *errors)
	return ae
}

// WithErrors sets the errors field of the AppError with a list of FieldError. This method allows for chaining and returns the modified AppError instance, enabling a fluent interface for building error responses.
func (ae *AppError) WithErrors(errors []FieldError) *AppError {
	ae.errors = errors
	return ae
}

// MarshalJSON implements the json.Marshaler interface for AppError. This method is responsible for converting the AppError instance into a JSON representation. It uses the standard library's json.Marshal function to serialize the AppError struct, ensuring that the resulting JSON includes the code and any errors if they are present.
func (ae *AppError) MarshalJSON() ([]byte, error) {
	type appErrorJSON struct {
		Errors []FieldError `json:"errors,omitempty"`
	}

	return json.Marshal(appErrorJSON{
		Errors: ae.errors,
	})
}

// WithField sets the field of the FieldError instance. This method allows for chaining and returns the modified FieldError instance, enabling a fluent interface for building error details.
func (e *FieldError) WithField(field string) *FieldError {
	e.field = field
	return e
}

func (e FieldError) MarshalJSON() ([]byte, error) {
	type errJSON struct {
		Field   string `json:"field,omitempty"`
		Message string `json:"message,omitempty"`
	}

	return json.Marshal(errJSON{
		Field:   e.field,
		Message: e.message,
	})
}


func (ae *AppError) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ae.GetCode())
	j, err := ae.MarshalJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(j)
}