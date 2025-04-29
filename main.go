package main

import (
	"log/slog"
	"net/http"

	"os"

	"github.com/TheCodeboy12/internal/server/handlers"

	"github.com/TheCodeboy12/internal/server/middlewere"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if groups != nil {
					return a
				}
				if a.Key == slog.MessageKey {
					a.Key = "message"
				} else if a.Key == slog.SourceKey {
					a.Key = "logging.googleapis.com/sourceLocation"
				} else if a.Key == slog.LevelKey {
					a.Key = "severity"
					level := a.Value.Any().(slog.Level)
					if level == 12 {
						a.Value = slog.StringValue("CRITICAL")
					}
				}
				return a
			},
		}))
	slog.SetDefault(logger)
}

func main() {
	bambooSecret := os.Getenv("BAMBOO_SECRET")
	if bambooSecret == "" {
		slog.Error("BAMBOO_SECRET is not set")
		return
	}
	envPort := os.Getenv("PORT")
	if envPort == "" {
		slog.Error("PORT is not set")
		return
	}

	router := http.NewServeMux()
	router.Handle("POST /", handlers.RootHandler())
	port := ":" + envPort

	srv := &http.Server{
		Addr:    port,
		Handler: middlewere.ValidateRequest(bambooSecret)(router),
	}
	slog.Info("Starting server", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err.Error())
		return
	}
}
