package api

import (
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Api interface {
	Mount(point string)
}

type jsonApi struct {
	cfg       utils.Config
	mux       *http.ServeMux
	subrouter *chi.Mux
	domain    domain.Domain
}

func NewApi(cfg utils.Config, mux *http.ServeMux, domain domain.Domain) *jsonApi {
	router := chi.NewMux()

	return &jsonApi{
		cfg:       cfg,
		mux:       mux,
		subrouter: router,
		domain:    domain,
	}
}

func (a *jsonApi) setupRoutes() {
	r := a.subrouter

	r.Post("/signup", a.postSignup)
}

func (a *jsonApi) Mount(point string) {
	a.setupRoutes()
	a.mux.Handle(point, a.subrouter)
}
