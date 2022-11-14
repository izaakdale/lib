package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var routes = []routeOptions{
	{http.MethodGet, "/ping", ping},
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

type routeOptions struct {
	method   string
	path     string
	function http.HandlerFunc
}

// NewRouter returns a http.Handler with the specified routesOptions.
// To be used in conjunction with WithRoute.
func New(opts ...routeOptions) http.Handler {
	router := httprouter.New()
	opts = append(opts, routes...)
	for _, route := range opts {
		router.HandlerFunc(route.method, route.path, route.function)
	}
	return router
}

// WithRoute takes a method and path string, as well as a HandlerFunc.
// Returns a routeOptions for inputting to the NewRouter function.
func WithRoute(m string, p string, f http.HandlerFunc) routeOptions {
	return routeOptions{m, p, f}
}
