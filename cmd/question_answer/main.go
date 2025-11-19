package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"question-answer/internal/config"
	"question-answer/internal/domain/qa"
	"question-answer/internal/infrastructure/http/handlers"
	mw "question-answer/internal/infrastructure/http/middleware"
	"question-answer/internal/infrastructure/storage/postgres"

	// "question-answer/internal/app"

	"question-answer/pkg/sl_logger/sl"
	"question-answer/pkg/sl_logger/slogpretty"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	log.Info("CHECKING DB Conn,", slog.String("Trying to connect with DSN", pgConfig.DSN))
	storage, err := postgres.New(pgConfig)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	_ = storage

	service := qa.NewService(storage)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(mw.NewMWLogger(log))

	r.Route("/questions", func(r chi.Router) {
		r.Get("/", handlers.NewGetQuestionHandler(log, service).ServeHTTP)
		r.Post("/", handlers.NewAddQuestionHandler(log, service).ServeHTTP)

		r.Route("/{questionID}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "questionID")
				handlers.NewGetAllQuestionHandler(log, service, id).ServeHTTP(w, r)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "questionID")
				handlers.NewDeleteQuestionHandler(log, service, id).ServeHTTP(w, r)
			})

			r.Route("/answers", func(r chi.Router) {
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					questionID := chi.URLParam(r, "questionID")
					handlers.NewAddAnswerHandler(log, service, questionID).ServeHTTP(w, r)
				})
			})
		})
	})

	// mux.Handle("/questions",
	// 	middleware.NewMWLogger(log)(
	// 		middleware.RequestID(
	// 			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 				switch r.Method {
	// 				case http.MethodGet:
	// 					handlers.NewGetQuestionHandler(log, service).ServeHTTP(w, r)
	// 				case http.MethodPost:
	// 					handlers.NewAddQuestionHandler(log, service).ServeHTTP(w, r)
	// 				default:
	// 					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 				}
	// 			}),
	// 		),
	// 	),
	// )

	// mux.Handle("/questions/",
	// 	middleware.NewMWLogger(log)(
	// 		middleware.RequestID(
	// 			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 				idStr := strings.TrimPrefix(r.URL.Path, "/questions/")
	// 				if idStr == "" || idStr == "/" {
	// 					http.Error(w, "id required", http.StatusBadRequest)
	// 					return
	// 				}
	// 				if r.Method == http.MethodGet {
	// 					handlers.NewGetAllQuestionHandler(log, service, idStr).ServeHTTP(w, r)
	// 					return
	// 				}
	// 				if r.Method == http.MethodDelete {
	// 					handlers.NewDeleteQuestionHandler(log, service, idStr).ServeHTTP(w, r)
	// 					return
	// 				}
	// 				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 			}),
	// 		),
	// 	),
	// )

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
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
