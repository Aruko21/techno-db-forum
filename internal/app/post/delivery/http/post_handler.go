package postHttp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/soulphazed/techno-db-forum/internal/app/general/respond"
	"github.com/soulphazed/techno-db-forum/internal/app/post"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"net/http"
)
type PostHandler struct {
	PostUsecase post.Usecase
}

func NewPostHandler(m *mux.Router, u post.Usecase) {
	handler := &PostHandler{
		PostUsecase: u,
	}

	m.HandleFunc("/post/{id}/details", handler.HandleUpdatePost).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/post/{id}/details", handler.HandleGetPostDetails).Methods(http.MethodGet)
}

func (h *PostHandler) HandleGetPostDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	postObj, err := h.PostUsecase.FindById(id, r.URL.Query())

	if err != nil || postObj == nil {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find post with id "+id+"\n"))
		return
	}

	respond.Respond(w, r, http.StatusOK, postObj)
}


func (h *PostHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	decoder := json.NewDecoder(r.Body)
	newPost := new(model.Post)
	err := decoder.Decode(newPost)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	postObj, err := h.PostUsecase.Update(id, newPost.Message)

	if err != nil || postObj == nil {
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find post with id "+id+"\n"))
		return
	}

	respond.Respond(w, r, http.StatusOK, postObj)
}
