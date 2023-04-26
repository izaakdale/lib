package router

import (
	"fmt"
	"net/http"
)

type (
	options struct {
		route      *routeOption
		middleware *middlewareOption
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

func checkMethodMiddleware(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf("%d %s not allowed\n", http.StatusMethodNotAllowed, r.Method)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// New returns a http.Handler with the specified routesOptions.
// To be used in conjunction with WithRoute and WithMiddleware.
func New(opts ...options) http.Handler {
	mux := http.NewServeMux()
	opts = append(opts, defaultOpts...)

	var middlewares []middlewareFunc

	for _, opt := range opts {
		if opt.route != nil {
			o := opt
			mux.Handle(o.route.path, checkMethodMiddleware(o.route.function, o.route.method))
		}
		if opt.middleware != nil {
			middlewares = append(middlewares, opt.middleware.function)
		}
	}

	if len(middlewares) > 0 {
		// need to kick off the handler func cascade with the router
		curr := middlewares[0](mux)
		// keep wrapping each of the remaining middleware functions around the current
		for i := 1; i < len(middlewares); i++ {
			curr = middlewares[i](curr)
		}
		return curr
	}
	return mux
}

// WithRoute takes a method and path string, as well as a HandlerFunc.
// Returns a routeOptions for inputting to the NewRouter function.
func WithRoute(m string, p string, f http.HandlerFunc) options {
	return options{route: &routeOption{m, p, f}}
}

// WithMiddleware adds a middleware function to the router for processing each request.
// When using this function multiple times the last entry will be called first with
// the rest in reverse order.
func WithMiddleware(mf middlewareFunc) options {
	return options{middleware: &middlewareOption{mf}}
}
