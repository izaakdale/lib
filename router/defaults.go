package router

import (
	"log"
	"net/http"
)

var defaultOpts = []routerOptions{
	{route: &routeOption{http.MethodGet, "/_/ping", ping}},
	{middleware: &middlewareOption{urlLogger}},
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
func urlLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hit: %+v\n", r.URL)
		next.ServeHTTP(w, r)
	})
}
