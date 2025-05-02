package middlewere

import (
	"log/slog"
	"net/http"
	"time"
)

// define new response writer
type wrapperWriter struct {
	http.ResponseWriter
	status int
}

// implement the write header function
func (w *wrapperWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		timeTaken := time.Since(start)
		wrapped := &wrapperWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		slog.Info(
			"Request received",
			"url", r.URL.Path,
			"time_taken", timeTaken,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"user_agent", r.UserAgent(),
			"response_status", wrapped.status,
		)
	})
}
