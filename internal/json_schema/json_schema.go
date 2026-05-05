package jsonschema

import (
	"net/http"

	"github.com/MXLange/go-model/internal/errors"
	"github.com/xeipuuv/gojsonschema"
)


type schemaValidation struct {
	schema string
}		


func New(schema string) *schemaValidation {
	return &schemaValidation{schema: schema}
}


func (s *schemaValidation) Validate(str string) *errors.AppError {
	
	schemaLoader := gojsonschema.NewStringLoader(s.schema)
	documentLoader := gojsonschema.NewStringLoader(str)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return errors.New(http.StatusInternalServerError).WithError(errors.NewError("", "failed to validate schema"))
	}

	if result.Valid() {
		return nil
	}

	var appErrs []errors.FieldError = make([]errors.FieldError, len(result.Errors()))
	for i, err := range result.Errors() {
		description := err.Description()
		field := err.Field()
		appErrs[i] = *errors.NewError(field, description)
	}

	return errors.New(http.StatusBadRequest).WithErrors(appErrs)
} 


