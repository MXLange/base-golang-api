package db

import (
	"context"
	"database/sql"

	"github.com/MXLange/go-model/internal/logger"
)

type App struct {
	name          string
	logger        logger.LoggerIF
	driverName    string
	connectionStr string
}

func NewApp(driverName, connectionStr string, logger logger.LoggerIF) (*App, error) {
	return &App{
		name:          "app",
		driverName:    driverName,
		connectionStr: connectionStr,
		logger:        logger,
	}, nil
}

func (a *App) GetName() string {
	return a.name
}

func (a *App) GetConnectionStr() string {
	return a.connectionStr
}

func (a *App) Connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open(a.driverName, a.connectionStr)
	if err != nil {
		a.logger.Errorf(ctx, "%s db failed to connect to the database", a.name)
		return nil, err
	}
	a.logger.Infof(ctx, "%s db successfully connected to the database", a.name)
	return db, nil
}
