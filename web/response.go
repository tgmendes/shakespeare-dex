package web

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a json response containing an error message with the reason.
type ErrorResponse struct {
	Message string `json:"error"`
}

func respond(w http.ResponseWriter, data []byte, statusCode int, contentType string) {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RespondJSON is a utility to convert a Go value to JSON and send it to the client.
func RespondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respond(w, jsonData, statusCode, "application/json")
}

// RespondError is a utility to create an error response and send it to the client.
func RespondError(w http.ResponseWriter, errMsg string, statusCode int) {
	errResp := ErrorResponse{errMsg}

	RespondJSON(w, errResp, statusCode)
}
