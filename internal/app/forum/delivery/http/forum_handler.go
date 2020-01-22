package forumHttp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/soulphazed/techno-db-forum/internal/app/general/respond"
	"github.com/soulphazed/techno-db-forum/internal/app/forum"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"net/http"
)

type ForumHandler struct {
	ForumUsecase forum.Usecase
}

func NewForumHandler(m *mux.Router, fu forum.Usecase) {
	handler := &ForumHandler{
		ForumUsecase: fu,
	}

	m.HandleFunc("/forum/create", handler.HandleForumCreate).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/forum/{slug}/create", handler.HandleForumCreateThread).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/forum/{slug}/details", handler.HandleForumGetDetails).Methods(http.MethodGet)
	m.HandleFunc("/forum/{slug}/threads", handler.HandleForumGetThreads).Methods(http.MethodGet)
	m.HandleFunc("/forum/{slug}/users", handler.HandleForumGetUsers).Methods(http.MethodGet)
}

func (h *ForumHandler) HandleForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleForumCreate<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	newForum := new(model.Forum)
	err := decoder.Decode(newForum)
	if err != nil {
		err = errors.Wrapf(err, "HandleForumCreate<-Decode")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	forumObj, code, err := h.ForumUsecase.CreateForum(newForum)

	if code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find user: " + newForum.Author))
		return
	}

	if code == http.StatusConflict {
		respond.Respond(w, r, http.StatusConflict, forumObj)
		return
	}

	respond.Respond(w, r, http.StatusCreated, forumObj)
}

func (h *ForumHandler) HandleForumCreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleForumCreateThread<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	vars := mux.Vars(r)
	slug := vars["slug"]

	decoder := json.NewDecoder(r.Body)
	newThread := new(model.NewThread)
	err := decoder.Decode(newThread)
	if err != nil {
		err = errors.Wrapf(err, "HandleForumCreateThread<-Decode")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	threadObj, code, err := h.ForumUsecase.CreateThread(slug, newThread)

	if code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	if code == http.StatusConflict {
		respond.Respond(w, r, http.StatusConflict, threadObj)
		return
	}

	respond.Respond(w, r, http.StatusCreated, threadObj)
}

func (h *ForumHandler) HandleForumGetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	forumObj, err := h.ForumUsecase.Find(slug)

	if err != nil || forumObj == nil {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find forum with slug "+slug+"\n"))
		return
	}

	respond.Respond(w, r, http.StatusOK, forumObj)
}

func (h *ForumHandler) HandleForumGetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	forumSlug := vars["slug"]

	threads, code, err := h.ForumUsecase.GetThreadsByForum(forumSlug, r.URL.Query())

	if err != nil || code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find threads for forum slug "+forumSlug+"\n"))
		return
	}

	respond.Respond(w, r, http.StatusOK, threads)
}

func (h *ForumHandler) HandleForumGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	forumSlug := vars["slug"]

	users, code, err := h.ForumUsecase.GetUsersByForum(forumSlug, r.URL.Query())

	if err != nil || users == nil || code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find users for forum slug "+forumSlug+"\n"))
		return
	}

	respond.Respond(w, r, http.StatusOK, users)
}
