package router

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var defaultOpts = []routerOptions{
	{Route: &routeOption{http.MethodGet, "/_/ping", ping}},
	{Middleware: &middlewareOption{urlLogger}},
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

type (
	routerOptions struct {
		Route      *routeOption
		Middleware *middlewareOption
	}
	routeOption struct {
		method   string
		path     string
		function http.HandlerFunc
	}

	middlewareFunc   = func(next http.Handler) http.Handler
	middlewareOption struct {
		function middlewareFunc
	}
)

// NewRouter returns a http.Handler with the specified routesOptions.
// To be used in conjunction with WithRoute.
func New(opts ...routerOptions) http.Handler {
	router := httprouter.New()
	opts = append(opts, defaultOpts...)

	var middlewares []middlewareFunc
	for _, opt := range opts {
		if opt.Route != nil {
			router.HandlerFunc(opt.Route.method, opt.Route.path, opt.Route.function)
		}
		if opt.Middleware != nil {
			middlewares = append(middlewares, opt.Middleware.function)
		}
	}

	if len(middlewares) > 0 {
		// need to kick off the handler func cascade with the router
		curr := middlewares[0](router)
		// keep wrapping each of the remaining middleware functions around the current
		for i := 1; i < len(middlewares); i++ {
			curr = middlewares[i](curr)
		}
		return curr
	}
	return router
}

// WithRoute takes a method and path string, as well as a HandlerFunc.
// Returns a routeOptions for inputting to the NewRouter function.
func WithRoute(m string, p string, f http.HandlerFunc) routerOptions {
	return routerOptions{Route: &routeOption{m, p, f}}
}

// WithMiddleware adds a middleware function to the router for processing each request.
// When using this function multiple times the last entry will be called first with
// the rest in decsending order before finally reaching the default URL logger.
func WithMiddleware(mf middlewareFunc) routerOptions {
	return routerOptions{Middleware: &middlewareOption{mf}}
}
