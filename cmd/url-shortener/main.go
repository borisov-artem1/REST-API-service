package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoadConfig()

	//_ = cfg
	// для сервера можно использовать агрегатор логов, который будет принимать json
	log := setupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))
	log.Info("starting server")
	log.Debug("this is debug log")

	storage, err := postgres.NewStorage(cfg.StoragePath)

	if err != nil {
		log.Error("could not connect to storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	err = storage.DeleteURL("https://hyperpc.com")

	if err != nil {
		log.Error("could not connect to storage", sl.Err(err))
	}

	log.Info("delete alias")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
