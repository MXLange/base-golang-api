package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	usermocks "github.com/MXLange/go-model/internal/domain/user/mocks"
	internalerrors "github.com/MXLange/go-model/internal/errors"
	loggermocks "github.com/MXLange/go-model/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewHandlers(t *testing.T) {
	t.Run("returns error when service is nil", func(t *testing.T) {
		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, err := NewHandlers("user", nil, loggerMock)

		require.Error(t, err)
		assert.ErrorIs(t, err, internalerrors.ErrNilService)
		assert.Nil(t, got)
	})

	t.Run("returns error when logger is nil", func(t *testing.T) {
		serviceMock := usermocks.NewMockServicesIF(t)

		got, err := NewHandlers("user", serviceMock, nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, internalerrors.ErrNilLogger)
		assert.Nil(t, got)
	})

	t.Run("returns handlers when params are valid", func(t *testing.T) {
		serviceMock := usermocks.NewMockServicesIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, err := NewHandlers("user", serviceMock, loggerMock)

		require.NoError(t, err)
		require.NotNil(t, got)

		handler, ok := got.(*handlers)
		require.True(t, ok)
		assert.Equal(t, "user", handler.name)
		assert.Equal(t, serviceMock, handler.service)
		assert.Equal(t, loggerMock, handler.logger)
	})
}

func TestHandlersGetName(t *testing.T) {
	handler := &handlers{name: "user"}

	assert.Equal(t, "user", handler.GetName())
}

func TestHandlersCreate(t *testing.T) {
	t.Run("returns bad request when payload does not satisfy schema", func(t *testing.T) {
		loggerMock := loggermocks.NewMockLoggerIF(t)
		serviceMock := usermocks.NewMockServicesIF(t)

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s handler - received request to create.", mock.Anything).
			Return()
		loggerMock.EXPECT().
			Errorf(mock.Anything, "%s handler - request body failed schema validation: %v", mock.Anything, mock.Anything).
			Return()

		handler := &handlers{
			name:    "user",
			logger:  loggerMock,
			service: serviceMock,
		}

		req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{}`))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "name is required")
		assert.Contains(t, rec.Body.String(), "nestedField is required")
	})

	t.Run("returns success when payload satisfies schema", func(t *testing.T) {
		loggerMock := loggermocks.NewMockLoggerIF(t)
		serviceMock := usermocks.NewMockServicesIF(t)

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s handler - received request to create.", mock.Anything).
			Return()

		handler := &handlers{
			name:    "user",
			logger:  loggerMock,
			service: serviceMock,
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/create",
			strings.NewReader(`{"name":"My User","nestedField":{"data":"Some data"}}`),
		)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "hello from user handler", rec.Body.String())
	})
}
