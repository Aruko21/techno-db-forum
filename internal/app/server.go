package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	forumHttp "github.com/soulphazed/techno-db-forum/internal/app/forum/delivery/http"
	forumRepository "github.com/soulphazed/techno-db-forum/internal/app/forum/repository"
	forumUsecase "github.com/soulphazed/techno-db-forum/internal/app/forum/usecase"
	postHttp "github.com/soulphazed/techno-db-forum/internal/app/post/delivery/http"
	postRepository "github.com/soulphazed/techno-db-forum/internal/app/post/repository"
	postUsecase "github.com/soulphazed/techno-db-forum/internal/app/post/usecase"
	serviceHttp "github.com/soulphazed/techno-db-forum/internal/app/service/delivery/http"
	serviceRepository "github.com/soulphazed/techno-db-forum/internal/app/service/repository"
	serviceUsecase "github.com/soulphazed/techno-db-forum/internal/app/service/usecase"
	threadHttp "github.com/soulphazed/techno-db-forum/internal/app/thread/delivery/http"
	threadRepository "github.com/soulphazed/techno-db-forum/internal/app/thread/repository"
	threadUsecase "github.com/soulphazed/techno-db-forum/internal/app/thread/usecase"
	userHttp "github.com/soulphazed/techno-db-forum/internal/app/user/delivery/http"
	userRepository "github.com/soulphazed/techno-db-forum/internal/app/user/repository"
	userUsecase "github.com/soulphazed/techno-db-forum/internal/app/user/usecase"
	"net/http"
)

type Server struct {
	Mux *mux.Router
	Config *Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func NewServer(config *Config) (*Server, error) {
	server := &Server {
		Mux: mux.NewRouter().PathPrefix("/api").Subrouter(),
		Config: config,
	}

	return server, nil
}

func (s *Server) ConfigureServer(db *pgx.ConnPool) {
	userRep := userRepository.NewUserRepository(db)
	userUse := userUsecase.NewForumUsecase(userRep)
	userHttp.NewUserHandler(s.Mux, userUse)

	threadRep := threadRepository.NewThreadRepository(db)
	threadUse := threadUsecase.NewThreadUsecase(threadRep, userRep)
	threadHttp.NewThreadHandler(s.Mux, threadUse)

	forumRep := forumRepository.NewForumRepository(db)
	forumUse := forumUsecase.NewForumUsecase(forumRep, userRep, threadRep)
	forumHttp.NewForumHandler(s.Mux, forumUse)

	postRep := postRepository.NewPostRepository(db)
	postUse := postUsecase.NewPostUsecase(postRep)
	postHttp.NewPostHandler(s.Mux, postUse)

	serviceRep := serviceRepository.NewServiceRepository(db)
	serviceUse := serviceUsecase.NewServiceUsecase(serviceRep)
	serviceHttp.NewServiceHandler(s.Mux, serviceUse)
}
