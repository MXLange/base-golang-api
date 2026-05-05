package app

import (
	"context"
	stderrors "errors"
	"testing"

	appmocks "github.com/MXLange/go-model/internal/domain/app/mocks"
	internalerrors "github.com/MXLange/go-model/internal/errors"
	loggermocks "github.com/MXLange/go-model/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewServices(t *testing.T) {
	t.Run("returns error when repository is nil", func(t *testing.T) {
		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, err := NewServices("app", nil, loggerMock)

		require.Error(t, err)
		assert.ErrorIs(t, err, internalerrors.ErrNilRepository)
		assert.Nil(t, got)
	})

	t.Run("returns error when logger is nil", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)

		got, err := NewServices("app", repositoryMock, nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, internalerrors.ErrNilLogger)
		assert.Nil(t, got)
	})

	t.Run("returns services when params are valid", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, err := NewServices("app", repositoryMock, loggerMock)

		require.NoError(t, err)
		require.NotNil(t, got)

		service, ok := got.(*services)
		require.True(t, ok)
		assert.Equal(t, "app", service.name)
		assert.Equal(t, repositoryMock, service.repository)
		assert.Equal(t, loggerMock, service.logger)
	})
}

func TestServicesHealth(t *testing.T) {
	t.Run("returns nil when repository ping succeeds", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s service - received request to ping the database.", mock.Anything).
			Return()
		repositoryMock.EXPECT().
			Ping(mock.Anything).
			Return(nil)

		service := &services{
			name:       "app",
			logger:     loggerMock,
			repository: repositoryMock,
		}

		err := service.Health(context.Background())

		require.NoError(t, err)
	})

	t.Run("returns repository ping error", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)
		expectedErr := stderrors.New("ping failed")

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s service - received request to ping the database.", mock.Anything).
			Return()
		repositoryMock.EXPECT().
			Ping(mock.Anything).
			Return(expectedErr)

		service := &services{
			name:       "app",
			logger:     loggerMock,
			repository: repositoryMock,
		}

		err := service.Health(context.Background())

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestServicesCreate(t *testing.T) {
	t.Run("returns repository result", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s service - received request to create.", mock.Anything).
			Return()
		repositoryMock.EXPECT().
			Create(mock.Anything, "murillo").
			Return(7, nil)

		service := &services{
			name:       "app",
			logger:     loggerMock,
			repository: repositoryMock,
		}

		gotID, gotErr := service.Create(context.Background(), "murillo")

		assert.Equal(t, 7, gotID)
		assert.Nil(t, gotErr)
	})

	t.Run("returns repository app error", func(t *testing.T) {
		repositoryMock := appmocks.NewMockRepositoryIF(t)
		loggerMock := loggermocks.NewMockLoggerIF(t)
		expectedErr := internalerrors.New(400).WithError(internalerrors.NewError("", "invalid name"))

		loggerMock.EXPECT().
			Infof(mock.Anything, "%s service - received request to create.", mock.Anything).
			Return()
		repositoryMock.EXPECT().
			Create(mock.Anything, "murillo").
			Return(0, expectedErr)

		service := &services{
			name:       "app",
			logger:     loggerMock,
			repository: repositoryMock,
		}

		gotID, gotErr := service.Create(context.Background(), "murillo")

		assert.Zero(t, gotID)
		assert.Same(t, expectedErr, gotErr)
	})
}
