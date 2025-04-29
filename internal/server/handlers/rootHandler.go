package handlers

import (
	"net/http"
)

func RootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// will post to the pubsub I guess

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
