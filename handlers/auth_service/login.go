package authservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
	"github.com/clevrf0x/go-multitenancy/server"
)

func LoginUserHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			EmailAddress string `json:"email_address"`
			Password     string `json:"password"`
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

		// get tenant id user present in
		s.DB.SwitchPublicSchema(r.Context())
		lookup, err := s.Queries.GetUserLookupByEmail(r.Context(), requestBody.EmailAddress)
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("user does not exist"),
				http.StatusNotFound,
			)
			return
		}

		// switch to tenant schema for checking password
		s.DB.SwitchSchema(r.Context(), lookup.TenantID.String())
		defer s.DB.SwitchPublicSchema(r.Context())

		user, err := s.Queries.GetUserByEmail(r.Context(), requestBody.EmailAddress)
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("user does not exist"),
				http.StatusNotFound,
			)
			return
		}

		if err := helpers.ComparePassword(user.Password.String, requestBody.Password); err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("invalid password"),
				http.StatusUnauthorized,
			)
			return
		}

		// issue tokens
		accessToken, err := helpers.GenerateAccessToken(user.ID.String(), lookup.TenantID.String())
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not generate access_token"),
				http.StatusInternalServerError,
			)
			return
		}
		refreshToken, err := helpers.GenerateRefreshToken(user.ID.String(), lookup.TenantID.String())
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not generate refresh_token"),
				http.StatusInternalServerError,
			)
			return
		}

		response := map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"tenant_id":     lookup.TenantID.String(),
		}
		helpers.WriteJSONResponse(w, helpers.NewAPISuccessResponse(response), http.StatusOK)
	}
}
