package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kudras3r/url_shortener/internal/config"
	"github.com/kudras3r/url_shortener/internal/http-server/handlers/redirect"
	"github.com/kudras3r/url_shortener/internal/http-server/handlers/save"
	"github.com/kudras3r/url_shortener/internal/lib/logger/sl"
	"github.com/kudras3r/url_shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	config := config.MustLoad()

	fmt.Println(config)

	logger := setupLogger(config.Env)
	logger.Debug("debug messages are enabled")

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		logger.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// router.Options("/url/", save.OptionsHandler)
	router.Post("/url/", save.New(logger, storage))
	router.Get("/url/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.Any("address", config.Address))

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Info("server stoped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return logger
}
