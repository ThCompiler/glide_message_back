package system_languages_handler

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app/delivery/http/models"
	usecase_info "glide/internal/app/usecase/info"
	bh "glide/internal/pkg/handler"
	"net/http"
)

type SystemLanguagesHandler struct {
	infoUsecase usecase_info.Usecase
	bh.BaseHandler
}

func NewSystemLanguagesHandler(log *logrus.Logger, ucInfo usecase_info.Usecase) *SystemLanguagesHandler {
	h := &SystemLanguagesHandler{
		BaseHandler: *bh.NewBaseHandler(log),
		infoUsecase: ucInfo,
	}

	h.AddMethod(http.MethodGet, h.GET)

	return h
}

func (s *SystemLanguagesHandler) GET(w http.ResponseWriter, r *http.Request) {
	info, err := s.infoUsecase.GetLanguages()
	if err != nil {
		s.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	s.Log(r).Debug("get info countryes")
	s.Respond(w, r, http.StatusOK, http_models.LanguagesToInfos(info))
}
