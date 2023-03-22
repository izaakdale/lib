package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetPathParam(r *http.Request, paramName string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(paramName)
}

func GetPathParams(r *http.Request, paramNames []string) map[string]string {
	params := httprouter.ParamsFromContext(r.Context())
	ret := make(map[string]string, len(paramNames))
	for _, n := range paramNames {
		param := params.ByName(n)
		ret[n] = param
	}
	return ret
}

func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func GetHeaders(r *http.Request, keys []string) map[string]string {
	headers := r.Header
	ret := make(map[string]string, len(keys))
	for _, k := range keys {
		h := headers.Get(k)
		ret[k] = h
	}
	return ret
}
