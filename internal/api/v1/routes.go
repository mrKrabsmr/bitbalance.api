package apiV1

import (
	"fl/my-portfolio/internal/app/controllers"

	"github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func ConfigureRoutes(controller *controllers.Controller, route chi.Router) {
	
	route.Get("/docs/*", httpSwagger.Handler(func(c *httpSwagger.Config) {
		httpSwagger.URL("http://0.0.0.0:8000/docs/swagger.json")
	}))

	route.Group(
		func(rr chi.Router) {
			rr.Use(middleware.Logger)
			rr.Use(middleware.Recoverer)
			rr.Post("/register", controller.Register)
			rr.Post("/login", controller.Login)
			rr.Post("/refresh", controller.Refresh)
			rr.Group(func(r chi.Router) {
				r.Use(controller.AuthenticationMiddleware)
				r.Get("/cryptocurrencies", controller.Cryptocurrencies)
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
