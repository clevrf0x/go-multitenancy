package middlewares

import (
	"context"
	"net/http"

	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
	"github.com/go-chi/jwtauth/v5"
)

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				helpers.WriteJSONResponse(
					w,
					helpers.NewAPIErrorResponse(err.Error()),
					http.StatusUnauthorized,
				)
				return
			}

			// Check token type
			if value, exists := claims["type"]; !exists || value != "access" {
				helpers.WriteJSONResponse(
					w,
					helpers.NewAPIErrorResponse("invalid token type"),
					http.StatusUnauthorized,
				)
				return
			}

			// Extract tenant_id claim
			tenantID, ok := claims["tenant_id"].(string)
			if !ok || tenantID == "" {
				helpers.WriteJSONResponse(
					w,
					helpers.NewAPIErrorResponse("tenant_id not found in token"),
					http.StatusUnauthorized,
				)
				return
			}

			// Store tenant_id in the context
			ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
