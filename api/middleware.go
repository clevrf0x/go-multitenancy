package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func RegisterMiddlewares(r *chi.Mux) {
	r.Use(cors.Handler(getCORSOptions()))
	r.Use(middleware.Logger, middleware.Recoverer, middleware.CleanPath)
}

func getCORSOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}
}
