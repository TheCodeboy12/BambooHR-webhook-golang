package middlewere

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/TheCodeboy12/internal/helpers"
)

const (
	maxTimeSkew = 3 * time.Minute
)

// validate the request before proceeding with expensive crypto operations
func ValidateRequest(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeStamp, ok := r.Header[http.CanonicalHeaderKey("x-bambooHR-timestamp")]
			if !ok {
				http.Error(w, "Missing X-BambooHR-Timestampp header", http.StatusBadRequest)
				return
			}
			signature, ok := r.Header[http.CanonicalHeaderKey("x-bamboohr-signature")]
			if !ok {
				http.Error(w, "Missing X-BambooHR-Signature header", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading body", http.StatusBadRequest)
				return
			}
			slog.Debug("Request body", "body", string(body))
			if !helpers.ValidateTimeStamp(timeStamp[0], float64(maxTimeSkew)) {
				http.Error(w, "Invalid timestamp", http.StatusBadRequest)
				return
			}
			payload := string(body) + timeStamp[0]
			if !helpers.ValidateHmac(payload, signature[0], secret) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
