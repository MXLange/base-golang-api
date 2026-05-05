package app

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	internalerrors "github.com/MXLange/go-model/internal/errors"
	loggermocks "github.com/MXLange/go-model/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NewRepository(t *testing.T) {
	t.Run("returns error when db is nil", func(t *testing.T) {
		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, err := NewRepository("app", nil, loggerMock)

		require.Error(t, err)
		assert.ErrorIs(t, err, internalerrors.ErrNilDB)
		assert.Nil(t, got)
	})

	t.Run("returns error when logger is nil", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		got, repoErr := NewRepository("app", db, nil)

		require.Error(t, repoErr)
		assert.ErrorIs(t, repoErr, internalerrors.ErrNilLogger)
		assert.Nil(t, got)
	})

	t.Run("returns repository when params are valid", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		loggerMock := loggermocks.NewMockLoggerIF(t)

		got, repoErr := NewRepository("app", db, loggerMock)

		require.NoError(t, repoErr)
		require.NotNil(t, got)

		repo, ok := got.(*repository)
		require.True(t, ok)
		assert.Equal(t, "app", repo.name)
		assert.Equal(t, db, repo.db)
		assert.Equal(t, loggerMock, repo.logger)
	})
}

func TestRepositoryPing(t *testing.T) {
	t.Run("returns nil when db ping succeeds", func(t *testing.T) {
		db, sqlDBMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		sqlDBMock.ExpectPing()

		loggerMock := loggermocks.NewMockLoggerIF(t)
		loggerMock.EXPECT().
			Infof(mock.Anything, "%s repository - received request to ping the database.", mock.Anything).
			Return()

		repo := &repository{
			name:   "app",
			logger: loggerMock,
			db:     db,
		}

		err = repo.Ping(context.Background())

		require.NoError(t, err)
		require.NoError(t, sqlDBMock.ExpectationsWereMet())
	})

	t.Run("returns ping error when db ping fails", func(t *testing.T) {
		db, sqlDBMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		expectedErr := stderrors.New("ping failed")
		sqlDBMock.ExpectPing().WillReturnError(expectedErr)

		loggerMock := loggermocks.NewMockLoggerIF(t)
		loggerMock.EXPECT().
			Infof(mock.Anything, "%s repository - received request to ping the database.", mock.Anything).
			Return()

		repo := &repository{
			name:   "app",
			logger: loggerMock,
			db:     db,
		}

		err = repo.Ping(context.Background())

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		require.NoError(t, sqlDBMock.ExpectationsWereMet())
	})
}

func TestRepositoryCreate(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	loggerMock := loggermocks.NewMockLoggerIF(t)
	loggerMock.EXPECT().
		Infof(mock.Anything, "%s repository - received request to create.", mock.Anything).
		Return()

	repo := &repository{
		name:   "app",
		logger: loggerMock,
		db:     db,
	}

	gotID, gotErr := repo.Create(context.Background(), "murillo")

	assert.Zero(t, gotID)
	assert.Nil(t, gotErr)
}
