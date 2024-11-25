package tenantservice

import (
	"net/http"

	"github.com/clevrf0x/go-multitenancy/pkg/helpers"
	"github.com/clevrf0x/go-multitenancy/server"
	"github.com/google/uuid"
)

func GetCurrentTenantHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, err := uuid.Parse(r.Context().Value("tenant_id").(string))
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("malformed tenant id"),
				http.StatusForbidden,
			)
			return
		}

		// NOTE: If needed switch schema with
		// s.DB.SwitchSchema(r.Context(), tenantID.String)

		tenant, err := s.Queries.GetTenantByID(r.Context(), tenantID)
		if err != nil {
			helpers.WriteJSONResponse(
				w,
				helpers.NewAPIErrorResponse("could not find tenant info"),
				http.StatusNotFound,
			)
			return
		}
		helpers.WriteJSONResponse(w, helpers.NewAPISuccessResponse(tenant), http.StatusOK)
	}
}
