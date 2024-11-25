package authservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clevrf0x/go-multitenancy/db"
	db_sqlc "github.com/clevrf0x/go-multitenancy/db/sqlc"
	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
	"github.com/clevrf0x/go-multitenancy/server"
)

func RegisterUserHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Name            string `json:"name"`
			EmailAddress    string `json:"email_address"`
			Password        string `json:"password"`
			TenantName      string `json:"tenant_name"`
			TenantWorkSpace string `json:"tenant_workspace"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("could not decode request body: %v\n", err)
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("Invalid request body"),
				http.StatusBadRequest,
			)
			return
		}

		// Create a new tenant and apply migrations
		tenant, err := db.CreateTenant(s.DB.DB, requestBody.TenantName, requestBody.TenantWorkSpace)
		if err != nil {
			log.Print(err)
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not create new tenant"),
				http.StatusInternalServerError,
			)
			return
		}

		// Switch to tenant schema
		s.DB.SwitchSchema(r.Context(), tenant.ID.String())

		// Create new user inside tenant and a lookup entry in public
		password, _ := helpers.HashPassword(requestBody.Password)
		user, err := s.Queries.CreateUser(r.Context(), db_sqlc.CreateUserParams{
			Name:     requestBody.Name,
			Email:    requestBody.EmailAddress,
			Password: helpers.NewNullString(password),
		})
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not create user inside tenant"),
				http.StatusInternalServerError,
			)
			return
		}

		// Switch to public schema before creating a user lookup entry
		s.DB.SwitchPublicSchema(r.Context())
		_, err = s.Queries.CreateUserLookup(r.Context(), db_sqlc.CreateUserLookupParams{
			TenantID: tenant.ID,
			Email:    user.Email,
		})
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not create user lookup entry"),
				http.StatusInternalServerError,
			)
		}

		// Send user and tenant info as response
		response := map[string]any{
			"user": map[string]any{
				"id":             user.ID.String(),
				"name":           user.Name,
				"email_address":  user.Email,
				"created_at":     user.CreatedAt,
				"last_active_at": user.LastActiveAt,
			},
			"tenant": map[string]any{
				"id":         tenant.ID.String(),
				"name":       tenant.Name,
				"workspace":  tenant.Workspace,
				"created_at": tenant.CreatedAt,
			},
		}
		helpers.WriteJSONResponse(w, helpers.NewAPISuccessResponse(response), http.StatusCreated)
	}
}
