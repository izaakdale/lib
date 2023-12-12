package router

import (
	"net/http"
)

var defaultOpts = []options{
	{route: &routeOption{http.MethodGet, "/_/ping", ping}},
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}
