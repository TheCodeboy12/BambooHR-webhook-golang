package main

import (
	"log/slog"
	"net/http"

	"os"

	"github.com/TheCodeboy12/bambooWebhook/internal/server/handlers"
	"github.com/TheCodeboy12/bambooWebhook/internal/server/middlewere"
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
		slog.Error("Env variable missing.", "key", "BAMBOO_SECRET")
		os.Exit(1)
	}
	envPort := os.Getenv("PORT")
	if envPort == "" {
		slog.Error("Env variable missing PORT is not set")
		os.Exit(1)
	}
	cloudProjectId := os.Getenv("CLOUD_PROJECT_ID")
	if cloudProjectId == "" {
		slog.Error("Env variable missing: CLOUD_PROJECT_ID is not set")
		os.Exit(1)
	}
	topicName := os.Getenv("TOPIC_NAME")
	if topicName == "" {
		slog.Error("Env variable missing: TOPIC_NAME is not set")
		os.Exit(1)
	}
	router := http.NewServeMux()
	router.Handle("POST /{$}", handlers.RootHandler(cloudProjectId, topicName))

	port := ":" + envPort

	srv := &http.Server{
		Addr: port,
		Handler: middlewere.LoggingMiddleware(
			middlewere.ValidateRequest(bambooSecret)(router),
		),
	}
	slog.Info("Starting server", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err.Error())
		os.Exit(1)
	}
}
