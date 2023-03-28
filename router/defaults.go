package router

import (
	"log"
	"net/http"

	"github.com/izaakdale/lib/response"
)

var defaultOpts = []routerOptions{
	{route: &routeOption{http.MethodGet, "/_/ping", ping}},
	{middleware: &middlewareOption{urlLogger}},
}

func ping(w http.ResponseWriter, r *http.Request) {
	response.WriteJson(w, http.StatusOK, "pong")
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
func urlLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)
		select {
		case <-r.Context().Done():
			// recorder does not record changes with http.TimeoutHandler.
			// This hack will undoubtably produce bugs later down the line (when contexts are cancelled for other reasons).
			log.Printf("hit: %+v with code: %+v", r.URL, http.StatusServiceUnavailable)
		default:
			log.Printf("hit: %+v with code: %+v", r.URL, recorder.Status)
		}

	})
}
