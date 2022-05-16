package handler

import (
	"github.com/gorilla/mux"
	hf "glide/internal/pkg/handler/handler_interfaces"
	"glide/internal/pkg/utilits"
	"glide/internal/pkg/utilits/delivery"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	DELETE  = http.MethodDelete
	OPTIONS = http.MethodOptions
)

type BaseHandler struct {
	handlerMethods map[string]hf.HandlerFunc
	middlewares    []hf.HMiddlewareFunc
	HelpHandlers
}

func NewBaseHandler(log *logrus.Logger) *BaseHandler {
	h := &BaseHandler{handlerMethods: map[string]hf.HandlerFunc{}, middlewares: []hf.HMiddlewareFunc{},
		HelpHandlers: HelpHandlers{
			ErrorConvertor: delivery.ErrorConvertor{
				Responder: delivery.Responder{
					LogObject: utilits.NewLogObject(log),
				},
			},
		},
	}
	return h
}

func (h *BaseHandler) AddMiddleware(middleware ...hf.HMiddlewareFunc) {
	h.middlewares = append(h.middlewares, middleware...)
}

func (h *BaseHandler) AddMethod(method string, handlerMethod hf.HandlerFunc, middlewares ...hf.HFMiddlewareFunc) {
	h.handlerMethods[method] = h.applyHFMiddleware(handlerMethod, middlewares...)
}

func (h *BaseHandler) applyHFMiddleware(handlerMethod hf.HandlerFunc,
	middlewares ...hf.HFMiddlewareFunc) hf.HandlerFunc {
	resultHandlerMethod := handlerMethod
	for index := len(middlewares) - 1; index >= 0; index-- {
		resultHandlerMethod = middlewares[index](resultHandlerMethod)
	}
	return resultHandlerMethod
}

func (h *BaseHandler) applyMiddleware(handler hf.Handler) hf.Handler {
	resultHandler := handler
	for index := len(h.middlewares) - 1; index >= 0; index-- {
		resultHandler = h.middlewares[index](resultHandler)
	}
	return resultHandler
}

func (h *BaseHandler) getListMethods() []string {
	var useMethods []string
	for key := range h.handlerMethods {
		useMethods = append(useMethods, key)
	}
	useMethods = append(useMethods, http.MethodOptions)
	return useMethods
}

func (h *BaseHandler) add(handler http.Handler, route *mux.Route) {
	var methods []string
	for key := range h.handlerMethods {
		methods = append(methods, key)
	}
	route.Handler(handler).Methods(methods...)
}

func (h *BaseHandler) Connect(route *mux.Route) {
	h.add(h.applyMiddleware(h), route)
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.PrintRequest(w, r)
	ok := true
	var hndlr hf.HandlerFunc

	hndlr, ok = h.handlerMethods[r.Method]
	if ok {
		hndlr(w, r)
	} else {
		h.Log(r).Errorf("Unexpected http method: %s", w, r.Method)
		r.Header.Set("Allow", strings.Join(h.getListMethods(), ", "))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
