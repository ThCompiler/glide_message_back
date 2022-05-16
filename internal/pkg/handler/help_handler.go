package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"glide/internal/app"
	repFiles "glide/internal/app/repository/files"
	"glide/internal/pkg/handler/handler_errors"
	"glide/internal/pkg/utilits/delivery"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Pagination struct {
	Limit int64
	Desc  bool
	Since string
}

const (
	EmptyQuery   = -2
	DefaultLimit = 100
)

type Sanitizable interface {
	easyjson.Unmarshaler
	Sanitize(sanitizer bluemonday.Policy)
}

type HelpHandlers struct {
	delivery.ErrorConvertor
}

func (h *HelpHandlers) PrintRequest(w http.ResponseWriter, r *http.Request) {
	h.Log(r).Infof("Request: %s. With method: %s. From URL: %s", w, r.Method, r.Host+r.URL.Path)
}

// GetInt64FromParam HTTPErrors
//		Status 400 handler_errors.InvalidParameters
func (h *HelpHandlers) GetInt64FromParam(w http.ResponseWriter, r *http.Request, name string) (int64, int, error) {
	vars := mux.Vars(r)
	number := vars[name]
	numberInt, err := strconv.ParseInt(number, 10, 64)
	if number == "" || err != nil {
		h.Log(r).Infof("can't get parametrs %s, was got %v)", name, number)
		return app.InvalidInt, http.StatusBadRequest, handler_errors.InvalidParameters
	}
	return numberInt, app.InvalidInt, nil
}

// GetPaginationFromQuery Expected api param:
// 	Default value for limit - 100
//	Param since query any false "start number of values"
// 	Param limit query uint64 false "number values to return"
//	Param desc  query bool false "
// Errors:
// 	Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetPaginationFromQuery(w http.ResponseWriter, r *http.Request) (*Pagination, int, error) {
	limit, code, err := h.GetInt64FromQueries(w, r, "limit")
	if err != nil {
		return nil, code, err
	}

	if limit == EmptyQuery {
		limit = DefaultLimit
	}

	desc := h.GetBoolFromQueries(w, r, "desc")

	since, info := h.GetStringFromQueries(w, r, "since")
	if info == EmptyQuery {
		since = ""
	}
	return &Pagination{Since: since, Desc: desc, Limit: limit}, app.InvalidInt, nil
}

// GetInt64FromQueries HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetInt64FromQueries(w http.ResponseWriter, r *http.Request, name string) (int64, int, error) {
	number := r.URL.Query().Get(name)
	if number == "" {
		return EmptyQuery, app.InvalidInt, nil
	}

	numberInt, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return app.InvalidInt, http.StatusBadRequest, handler_errors.InvalidQueries
	}

	return numberInt, app.InvalidInt, nil
}

// GetBoolFromQueries HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetBoolFromQueries(w http.ResponseWriter, r *http.Request, name string) bool {
	number := r.URL.Query().Get(name)
	if number == "" {
		return false
	}

	numberInt, err := strconv.ParseBool(number)
	if err != nil {
		return false
	}

	return numberInt
}

// GetStringFromQueries HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetStringFromQueries(w http.ResponseWriter, r *http.Request, name string) (string, int) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return "", EmptyQuery
	}

	return value, app.InvalidInt
}

// GetStringFromParam HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetStringFromParam(w http.ResponseWriter, r *http.Request, name string) (string, int) {
	vars := mux.Vars(r)
	value := vars[name]
	if value == "" {
		return "", EmptyQuery
	}

	return value, app.InvalidInt
}

// GetArrayStringFromQueries HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetArrayStringFromQueries(w http.ResponseWriter, r *http.Request, name string) ([]string, int) {
	values := r.URL.Query().Get(name)
	if values == "" {
		return nil, EmptyQuery
	}

	return strings.Split(values, ","), app.InvalidInt
}

func (h *HelpHandlers) GetRequestBody(r *http.Request, reqStruct Sanitizable, sanitizer bluemonday.Policy) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if err := easyjson.UnmarshalFromReader(r.Body, reqStruct); err != nil {
		return err
	}

	reqStruct.Sanitize(sanitizer)
	return nil
}

// GetFilesFromRequest http Errors:
// 		Status 400 handler_errors.FileSizeError
// 		Status 400 handler_errors.InvalidFormFieldName
// 		Status 400 handler_errors.InvalidImageExt
// 		Status 500 handler_errors.InternalError
func (h *HelpHandlers) GetFilesFromRequest(w http.ResponseWriter, r *http.Request, maxSize int64,
	name string, validTypes []string) (io.Reader, repFiles.FileName, int, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, "", http.StatusBadRequest, app.GeneralError{
			ExternalErr: errors.Wrapf(err, "max size is : %d ", maxSize),
			Err:         handler_errors.FileSizeError,
		}
	}

	f, fHeader, err := r.FormFile(name)
	if err != nil {
		return nil, "", http.StatusBadRequest, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InvalidFormFieldName,
		}
	}

	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		return nil, "", http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		}
	}

	sort.Strings(validTypes)
	fType := http.DetectContentType(buff)
	if pos := sort.SearchStrings(validTypes, fType); pos == len(validTypes) || validTypes[pos] != fType {
		return nil, "", http.StatusBadRequest, fmt.Errorf("%s, %s",
			handler_errors.InvalidExt, strings.Join(validTypes, " ,"))
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return nil, "", http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		}
	}

	return f, repFiles.FileName(fHeader.Filename), 0, nil
}
