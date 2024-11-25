package handlers

import (
	"net/http"

	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
)

func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{"ping": "pong"}
		helpers.WriteJSONResponse(w, helpers.NewAPISuccessResponse(response), http.StatusOK)
	}
}

func ProtectedRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{"logged_in": true}
		helpers.WriteJSONResponse(
			w,
			helpers.NewAPISuccessResponse(response),
			http.StatusOK,
		)
	}
}
