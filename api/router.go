package api

import (
	"fmt"
	"net/http"

	"github.com/clevrf0x/go-multitenancy/api/middlewares"
	"github.com/clevrf0x/go-multitenancy/handlers"
	authservice "github.com/clevrf0x/go-multitenancy/handlers/auth_service"
	tenantservice "github.com/clevrf0x/go-multitenancy/handlers/tenant_service"
	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
	"github.com/clevrf0x/go-multitenancy/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func RegisterRoutes(s *server.Server) {
	router := s.Router

	tokenAuth, err := helpers.GetJWTAuth()
	if err != nil {
		panic(
			fmt.Errorf(
				"something wrong with jwt setup, check if environment variables set correct: %v",
				err.Error(),
			),
		)
	}

	// Setup custom error handlers
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSONResponse(
			w,
			helpers.NewAPIErrorResponse("route does not exist"),
			http.StatusNotFound,
		)
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSONResponse(
			w,
			helpers.NewAPIErrorResponse("method not allowed"),
			http.StatusMethodNotAllowed,
		)
	})

	// Create subrouter with "/api/v1" prefix
	r := chi.NewRouter()
	router.Mount("/api/v1", r)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/ping", handlers.HealthCheckHandler())

		r.Post("/register", authservice.RegisterUserHandler(s))
		r.Post("/login", authservice.LoginUserHandler(s))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middlewares.Authenticator(tokenAuth))

		r.Get("/protected", handlers.ProtectedRouteHandler())
		r.Get("/current_tenant", tenantservice.GetCurrentTenantHandler(s))
	})

}
