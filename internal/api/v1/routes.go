package apiV1

import (
	"fl/my-portfolio/internal/app/controllers"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func ConfigureRoutes(controller *controllers.Controller, route chi.Router) {
	route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	route.Get("/docs/*", httpSwagger.Handler(func(c *httpSwagger.Config) {
		httpSwagger.URL("http://0.0.0.0:8000/docs/swagger.json")
	}))

	route.Group(
		func(rr chi.Router) {
			rr.Use(middleware.Logger)
			rr.Use(middleware.Recoverer)

			fs := http.FileServer(http.Dir("static"))
			rr.Handle("/static/*", http.StripPrefix("/api/v1/static/", fs))

			rr.Get("/cryptocurrencies", controller.Cryptocurrencies)
			rr.Post("/register", controller.Register)
			rr.Post("/login", controller.Login)
			rr.Post("/refresh", controller.Refresh)
			rr.Group(func(r chi.Router) {
				r.Use(controller.AuthenticationMiddleware)
				r.Get("/getMe", controller.GetMe)
				r.Get("/portfolio", controller.GetPortfolio)
				r.Post("/portfolio", controller.CreatePortfolio)
				r.Patch("/portfolio/{id}", controller.UpdatePortfolioDetail)
				r.Delete("/portfolio/{id}", controller.DeletePortfolioDetail)
				r.Get("/stakings", controller.Stakings)
				r.Get("/stakings/portfolio", controller.StakingsPortfolio)
				r.Get("/stakings/detail/{crypt_symbol}", controller.StakingsDetail)
			})
		},
	)
}
