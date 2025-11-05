package routes

import (
	"github.com/diagnosis/luxsuv-api-v2/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetRouter(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(hr chi.Router) {
		hr.Get("/healthz", app.ServerHealthCheckerHandler.HandleHealth)
	})
	return r
}
