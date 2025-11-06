package routes

import (
	"github.com/diagnosis/luxsuv-api-v2/internal/app"
	"github.com/diagnosis/luxsuv-api-v2/internal/logger"
	customMiddleware "github.com/diagnosis/luxsuv-api-v2/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetRouter(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logger.HandlerLogger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", app.ServerHealthCheckerHandler.HandleHealth)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/login", app.UserHandler.HandleLogin)

		api.Group(func(protected chi.Router) {
			protected.Use(customMiddleware.RequireJWT(app.Signer))
		})

		api.Group(func(adminOnly chi.Router) {
			adminOnly.Use(customMiddleware.RequireJWT(app.Signer))
			adminOnly.Use(customMiddleware.RequireRole("admin"))
		})

		api.Group(func(driverOnly chi.Router) {
			driverOnly.Use(customMiddleware.RequireJWT(app.Signer))
			driverOnly.Use(customMiddleware.RequireRole("driver"))
		})
	})

	return r
}
