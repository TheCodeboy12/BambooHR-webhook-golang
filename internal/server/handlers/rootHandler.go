package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"cloud.google.com/go/pubsub"
)

func RootHandler(cloudProjectId string, topicName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		client, err := pubsub.NewClient(ctx, cloudProjectId)
		if err != nil {
			slog.Error("Error creating pubsub client", "error", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer client.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading body", "error", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		defer r.Body.Close()

		topic := client.Topic(topicName)
		//post message
		result := topic.Publish(ctx, &pubsub.Message{
			Data: body,
		})
		id, err := result.Get(ctx)
		if err != nil {
			slog.Error("Error publishing message", "error", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			slog.Info("Message published", "id", id)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
