package userHttp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/soulphazed/techno-db-forum/internal/app/general/respond"
	"github.com/soulphazed/techno-db-forum/internal/app/user"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"net/http"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHandler(m *mux.Router, u user.Usecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	m.HandleFunc("/user/{nickname}/create", handler.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleGetProfile).Methods(http.MethodGet)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleChangeProfile).Methods(http.MethodPost, http.MethodOptions)
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	newUser := new(model.User)
	err := decoder.Decode(newUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Decode:")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	newUser.Nickname = nickname

	users, err := h.UserUsecase.CreateUser(newUser)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateUser: ")
		respond.Error(w, r, http.StatusBadRequest, err)

		return
	}

	if users != nil {
		respond.Respond(w, r, http.StatusConflict, users)

		return
	}

	respond.Respond(w, r, http.StatusCreated, newUser)
}

func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	userObj, err := h.UserUsecase.Find(nickname)

	if err != nil || userObj == nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateUser: ")
		respond.Error(w, r, http.StatusNotFound, errors.New("Can't find user with nickname #'"+nickname+"'\n"))

		return
	}

	respond.Respond(w, r, http.StatusOK, userObj)
}

func (h *UserHandler) HandleChangeProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleChangeProfile<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	newUser := new(model.User)
	err := decoder.Decode(newUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleChangeProfile<-Decode:")
		respond.Error(w, r, http.StatusBadRequest, err)

		return
	}
	newUser.Nickname = nickname

	newUser, err, code := h.UserUsecase.Update(newUser)

	if code == http.StatusNotFound {
		respond.Error(w, r, http.StatusNotFound, err)

		return
	}

	if err != nil || code == http.StatusConflict {
		err = errors.Wrapf(err, "HandleChangeProfile<-Update: ")
		respond.Error(w, r, http.StatusConflict, err)

		return
	}

	respond.Respond(w, r, http.StatusOK, newUser)
}