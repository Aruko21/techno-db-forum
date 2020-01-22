package serviceHttp

import (
	"github.com/gorilla/mux"
	"github.com/soulphazed/techno-db-forum/internal/app/general/respond"
	"github.com/soulphazed/techno-db-forum/internal/app/service"
	"net/http"
)

type ServiceHandler struct {
	ServiceUsecase service.Usecase
}

func NewServiceHandler(m *mux.Router, u service.Usecase) {
	handler := &ServiceHandler{
		ServiceUsecase: u,
	}

	m.HandleFunc("/service/clear", handler.HandleServiceClear).Methods(http.MethodPost)
	m.HandleFunc("/service/status", handler.HandleServiceGetStatus).Methods(http.MethodGet)
}

func (h *ServiceHandler) HandleServiceClear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.ServiceUsecase.ClearAll()

	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, nil)
}

func (h *ServiceHandler) HandleServiceGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, err := h.ServiceUsecase.GetStatus()

	if err != nil || status == nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, status)
}
