package routes

import (
	"backend/internal/utils"
	"net/http"
)

func handleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, true, "API is healthy:)", nil)
}