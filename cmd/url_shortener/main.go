package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/kudras3r/url_shortener/internal/config"
	"github.com/kudras3r/url_shortener/internal/lib/logger/sl"
	"github.com/kudras3r/url_shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	// TODO: init config: cleanenv
	config := config.MustLoad()

	fmt.Println(config)

	// TODO: init logger
	logger := setupLogger(config.Env)
	logger.Debug("debug messages are enabled")

	// TODO: init storage
	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		logger.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}

	// TESTS:
	// err = storage.SaveURL("asd1.com", "asd1")
	// if err != nil {
	// 	logger.Error("failed to save url", sl.Err(err))
	// 	os.Exit(1)
	// }
	// logger.Info("saved url")
	// err = storage.SaveURL("asd.com", "asd")
	// if err != nil {
	// 	logger.Error("failed to save url", sl.Err(err))
	// 	os.Exit(1)
	// }

	_ = storage
	// TODO: init router

	// TODO: run server
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
