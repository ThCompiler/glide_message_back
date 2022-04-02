package handler_interfaces

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
type HMiddlewareFunc func(http.Handler) http.Handler
type HFMiddlewareFunc func(HandlerFunc) HandlerFunc
