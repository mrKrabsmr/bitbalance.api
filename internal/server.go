package server

import (
	"embed"
	apiV1 "fl/my-portfolio/internal/api/v1"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/controllers"
	"fl/my-portfolio/internal/configs"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed migrations
var EmbedMigrations embed.FS

type APIServer struct {
	logger *logrus.Logger
	config *configs.Config
	router *chi.Mux
}

func NewAPIServer() *APIServer {
	return &APIServer{
		config: core.GetConfig(),
		logger: core.GetLogger(),
		router: chi.NewRouter(),
	}
}

func (s *APIServer) configureRoutes() {
    controller := controllers.NewController()

	s.router.Route("/api", func(rr chi.Router) {

        rr.Route("/v1", func (r chi.Router) {
            apiV1.ConfigureRoutes(controller, r)
        })
	})
}

func (s *APIServer) StartMigrations() {
	conn := core.GetDB()

	goose.SetBaseFS(EmbedMigrations)
	if err := goose.SetDialect(s.config.DB.DBDialect); err != nil {
		panic(err)
	}

	if err := goose.Up(conn.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (s *APIServer) DownMigrations() {
	conn := core.GetDB()

    goose.SetBaseFS(EmbedMigrations)
	if err := goose.SetDialect(s.config.DB.DBDialect); err != nil {
		panic(err)
	}

	if err := goose.Down(conn.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (s *APIServer) Run() {
	s.configureRoutes()

	server := http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	s.logger.Infof("STARTING SERVER AT %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
