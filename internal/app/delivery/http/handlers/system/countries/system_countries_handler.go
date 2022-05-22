package system_countries_handler

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app/delivery/http/models"
	usecase_info "glide/internal/app/usecase/info"
	bh "glide/internal/pkg/handler"
	"net/http"
)

type SystemCountriesHandler struct {
	infoUsecase usecase_info.Usecase
	bh.BaseHandler
}

func NewSystemCountriesHandler(log *logrus.Logger, ucInfo usecase_info.Usecase) *SystemCountriesHandler {
	h := &SystemCountriesHandler{
		BaseHandler: *bh.NewBaseHandler(log),
		infoUsecase: ucInfo,
	}

	h.AddMethod(http.MethodGet, h.GET)

	return h
}

func (s *SystemCountriesHandler) GET(w http.ResponseWriter, r *http.Request) {
	info, err := s.infoUsecase.GetCountries()
	if err != nil {
		s.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	s.Log(r).Debug("get info countryes")
	s.Respond(w, r, http.StatusOK, http_models.CountriesToInfos(info))
}
