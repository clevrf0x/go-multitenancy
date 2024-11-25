package helpers

type APIResponse map[string]any

func NewAPISuccessResponse(data any) APIResponse {
	response := APIResponse{
		"status": "success",
	}
	if data != nil {
		response["data"] = data
	}
	return response
}

func NewAPIErrorResponse(reason any) APIResponse {
	response := APIResponse{
		"status": "error",
		"reason": reason,
	}
	return response
}
