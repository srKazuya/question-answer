package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"question-answer/internal/config"
	"question-answer/internal/infrastructure/storage/postgres"
	// "question-answer/internal/app"

	"question-answer/pkg/sl_logger/sl"
	"question-answer/pkg/sl_logger/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//entry point

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	pgConfig := postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s",
			cfg.DataBase.Host,
			cfg.DataBase.User,
			cfg.DataBase.Port,
			cfg.DataBase.Password,
			cfg.DataBase.Dbname,
			cfg.DataBase.Sslmode,
		),
		MigrationsPath: "internal/infrastructure/storage/postgres/migrations",
	}

	log.Info("adas,", slog.String("Trying to connect with DSN", pgConfig.DSN))
	storage, err := postgres.New(pgConfig)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	_ = storage

	// service := app.NewService(storage)
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting HTTP server",
		slog.String("address", cfg.Address),
		slog.Duration("read_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("write_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("idle_timeout", cfg.IdleTimeout),
	)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
