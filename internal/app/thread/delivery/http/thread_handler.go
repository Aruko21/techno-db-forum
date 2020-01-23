package threadHttp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/soulphazed/techno-db-forum/internal/app/general/respond"
	"github.com/soulphazed/techno-db-forum/internal/app/thread"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"net/http"
)

type ThreadHandler struct {
	ThreadUsecase thread.Usecase
}

func NewThreadHandler(m *mux.Router, u thread.Usecase) {
	handler := &ThreadHandler{
		ThreadUsecase: u,
	}

	m.HandleFunc("/thread/{slug_or_id}/create", handler.HandleCreatePosts).Methods(http.MethodPost)
	m.HandleFunc("/thread/{slug_or_id}/details", handler.HandleGetThreadDetails).Methods(http.MethodGet)
	m.HandleFunc("/thread/{slug_or_id}/details", handler.HandleUpdateThreadDetails).Methods(http.MethodPost)
	m.HandleFunc("/thread/{slug_or_id}/posts", handler.HandleGetThreadPosts).Methods(http.MethodGet)
	m.HandleFunc("/thread/{slug_or_id}/vote", handler.HandleVoteForThread).Methods(http.MethodPost)
}

func (h *ThreadHandler) HandleCreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	decoder := json.NewDecoder(r.Body)
	newPosts := new(model.Posts)
	err := decoder.Decode(newPosts)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	posts, code, err := h.ThreadUsecase.CreatePosts(slugOrId, newPosts)

	if code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	if code == http.StatusConflict {
		respond.Error(w, r, http.StatusConflict, err)
		return
	}

	respond.Respond(w, r, http.StatusCreated, posts)
}

func (h *ThreadHandler) HandleGetThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	threadObj, err := h.ThreadUsecase.FindByIdOrSlug(slugOrId)
	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, threadObj)
}

func (h *ThreadHandler) HandleUpdateThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	decoder := json.NewDecoder(r.Body)
	threadUpdate := new(model.ThreadUpdate)
	err := decoder.Decode(threadUpdate)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	threadObj, err := h.ThreadUsecase.UpdateThread(slugOrId, threadUpdate)
	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, threadObj)
}

func (h *ThreadHandler) HandleGetThreadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	posts, err := h.ThreadUsecase.GetThreadPosts(slugOrId, r.URL.Query())

	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, &posts)
}

func (h *ThreadHandler) HandleVoteForThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	decoder := json.NewDecoder(r.Body)
	vote := new(model.Vote)
	err := decoder.Decode(vote)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	threadObj, err := h.ThreadUsecase.Vote(slugOrId, vote)
	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, threadObj)
}