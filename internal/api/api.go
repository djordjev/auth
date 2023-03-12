package api

import (
	"github.com/djordjev/auth/internal/api/middleware"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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
	logger    *zap.SugaredLogger
}

func NewApi(cfg utils.Config, mux *http.ServeMux, domain domain.Domain, logger *zap.SugaredLogger) *jsonApi {
	router := chi.NewMux()

	return &jsonApi{
		cfg:       cfg,
		mux:       mux,
		subrouter: router,
		domain:    domain,
		logger:    logger,
	}
}

func (a *jsonApi) setupMiddleware() {
	r := a.subrouter

	r.Use(middleware.Logger(a.logger))

	r.Use(chimiddleware.Recoverer)
}

func (a *jsonApi) setupRoutes() {
	r := a.subrouter

	r.Post("/signup", a.postSignup)
}

func (a *jsonApi) Mount(point string) {
	a.setupMiddleware()
	a.setupRoutes()
	a.mux.Handle(point, a.subrouter)
}
