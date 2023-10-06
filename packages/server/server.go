package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/djordjev/auth/internal/api"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/notify"
	"github.com/djordjev/auth/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type server struct {
	mux    *http.ServeMux
	api    api.Api
	config utils.Config
	pool   *pgxpool.Pool
}

func (s *server) Mount(url string) {
	s.api.Mount(url)
}

func (s *server) Close() {
	s.pool.Close()
}

func NewServer(mux *http.ServeMux, config utils.Config) *server {
	srv := &server{mux: mux, config: config}

	srv.setup()

	return srv
}

func (s *server) setup() {
	// Init database
	dbUrl := s.config.GetConnectionString()
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		panic(err)
	}

	// Init redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", s.config.RedisHost, s.config.RedisPort),
		Password: s.config.RedisPassword,
		DB:       s.config.RedisDatabase,
	})

	// Setup repos
	repo := models.NewRepository(pool, client)

	// Init api
	logger := utils.MustBuildLogger(s.config)
	var notifier domain.Notifier
	if s.config.IsDev() {
		notifier = notify.SilentNotifier{}
	} else {
		notifier = notify.NewMailjetNotifier(s.config)
	}

	// Init app domain
	appDomain := domain.NewDomain(repo, s.config, notifier)

	// Init API
	appApi := api.NewApi(s.config, s.mux, appDomain, logger)
	s.api = appApi
}
