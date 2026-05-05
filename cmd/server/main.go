package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MXLange/go-model/env"
	"github.com/MXLange/go-model/internal/builder"
	"github.com/MXLange/go-model/internal/domain/app"
	"github.com/MXLange/go-model/internal/infra/db"
	"github.com/MXLange/go-model/internal/logger"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()
	l := logger.NewLogger()

	l.Info(ctx, "Starting server...")

	e, err := env.New(ctx, l)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	appDB, err := db.NewApp(e.DBDriverName, e.DBConnectionString, l)
	if err != nil {
		panic(err)
	}

	dbConn, err := appDB.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	builder, err := builder.New(r, l)
	if err != nil {
		panic(err)
	}

	app, err := app.NewApp(ctx, dbConn, l)
	if err != nil {
		panic(err)
	}

	builder.AddDomainEntity(app)

	if err := builder.Build(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := builder.Close(ctx); err != nil {
			l.Errorf(ctx, "Error closing builder: %v", err)
		}
	}()

	server := &http.Server{
		Addr:    ":" + e.Port,
		Handler: r,
	}

	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownCtx.Done()

		l.Info(ctx, "Shutdown signal received")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			l.Errorf(ctx, "server shutdown failed: %v\n", err)
		}
	}()

	l.Infof(ctx, "HTTP server listening on :%s", e.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	l.Info(ctx, "Server stopped")
}
