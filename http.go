package main

import (
	"encoding/json"
	"net/http"
)

// SuccessWithJSON sends a successful response with a JSON payload
func SuccessWithJSON(w http.ResponseWriter, payload interface{}) {
	RespondWithJSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"result": payload,
			"error":  nil,
		},
	)
}

// RespondWithJSON encodes the data as JSON, and sends a response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		RespondServerError(w, []string{"Failed to marshal JSON response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondServerError responds to the request with internal server error 500
func RespondServerError(w http.ResponseWriter, errorMessages []string) {
	RespondWithJSON(
		w,
		http.StatusInternalServerError,
		map[string][]string{
			"result": nil,
			"error":  errorMessages,
		},
	)
}
