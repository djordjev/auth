package api

import (
	"net/http"

	"github.com/djordjev/auth/internal/api/middleware"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type Api interface {
	Mount(point string)
}

type jsonApi struct {
	cfg       utils.Config
	mux       *http.ServeMux
	subrouter *chi.Mux
	domain    domain.Domain
	logger    *slog.Logger
}

func NewApi(cfg utils.Config, mux *http.ServeMux, domain domain.Domain, logger *slog.Logger) *jsonApi {
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
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
}

func (a *jsonApi) setupRoutes() {
	r := a.subrouter

	r.Post("/signup", a.postSignup)
	r.Post("/login", a.postLogin)
	r.Delete("/delete", a.deleteAccount)
	r.Post("/verify", a.postVerifyAccount)
	r.Post("/forget", a.postForgetPassword)
	r.Post("/passwordreset", a.postVerifyPasswordReset)
	r.Get("/session", a.getSession)
	r.Post("/logout", a.postLogout)
}

func (a *jsonApi) Mount(point string) {
	a.setupMiddleware()
	a.setupRoutes()
	a.mux.Handle(point, a.subrouter)
}
