package server

import "net/http"

type server struct {
	mux *http.ServeMux
}

func (s *server) Mount(url string) {
	s.mux.Handle(url, nil)
}

func NewServer(mux *http.ServeMux) *server {
	srv := &server{mux: mux}

	return srv
}
