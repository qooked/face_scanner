package main

import (
	"context"
	"faceScanner/config"
	"faceScanner/internal/adapters/repository"
	"faceScanner/internal/controllers/http"
	"faceScanner/internal/usecase"
	"faceScanner/pkg/database"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("config.LoadConfig(...): %w", err)
		return
	}

	postgres, err := database.Connect(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DatabaseName,
		cfg.Postgres.SslMode,
	)
	if err != nil {
		slog.Error("database.Connect(...): %w", err)
		return
	}
	slog.Info("Database connected")

	Server := http.NewServer(
		cfg.Server.Host,
		cfg.Server.Port,
		cfg.Server.AuthorizationKey,
	)

	repo := repository.New(postgres)
	uc := usecase.New(repo)
	Server.AttachHandlers(ctx, uc)

	doneChan := make(chan struct{})
	go Server.Run()
	slog.Info(fmt.Sprintf("Server started on %s:%s", cfg.Server.Host, cfg.Server.Port))
	go func() {
		<-signalChan
		slog.Info("Shutting down server...")
		err = Server.Shutdown(ctx)
		if err != nil {
			slog.Error("Server.Shutdown(...): %w", err)
		}
		doneChan <- struct{}{}
	}()
	<-doneChan
	cancel()
	slog.Info("Application stopped")

}
