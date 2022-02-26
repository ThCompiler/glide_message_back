package delivery

import (
	"github.com/mailru/easyjson"
	"glide/internal/pkg/utilits"
	"net/http"
)

//go:generate easyjson -disallow_unknown_fields responder.go

//easyjson:json
type ErrResponse struct {
	Err string `json:"message"`
}

type Responder struct {
	utilits.LogObject
}

func (h *Responder) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, ErrResponse{Err: err.Error()})
}

func (h *Responder) Respond(w http.ResponseWriter, r *http.Request, code int, data easyjson.Marshaler) {
	w.WriteHeader(code)
	if data != nil {
		_, _, err := easyjson.MarshalToHTTPResponseWriter(data, w)
		if err != nil {
			//h.Log(w, r).Error(jw.Error)
		}
	}
	//logUser, _ := easyjson.Marshal(data)
	//h.Log(w, r).Info("Respond data: ", string(logUser))
}
