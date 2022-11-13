package router

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var routes = []routeOptions{
	{http.MethodGet, "/ping", ping},
}

func ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit pingpong\n")
	w.Write([]byte("pong"))
}

type routeOptions struct {
	method   string
	path     string
	function http.HandlerFunc
}

func NewRouter(opts ...routeOptions) http.Handler {
	router := httprouter.New()
	opts = append(opts, routes...)
	for _, route := range opts {
		router.HandlerFunc(route.method, route.path, route.function)
	}
	return router
}

func WithRoute(m string, p string, f http.HandlerFunc) routeOptions {
	return routeOptions{m, p, f}
}
