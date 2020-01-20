package apiserver

import (
	"github.com/gorilla/mux"
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
		Mux: mux.NewRouter(),
		Config: config,
	}

	return server, nil
}

