package utils

import (
	"encoding/json"
	"net/http"
)


type APIResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
}	

func JSON(w http.ResponseWriter, status int, success bool, message string, data any) {
	resp:= APIResponse{
		Status: status,
		Success: success,
		Message: message,
		Data: data,
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}