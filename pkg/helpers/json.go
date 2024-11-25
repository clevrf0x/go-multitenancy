package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, response APIResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Error writing response body: %v", err)
	}
}
