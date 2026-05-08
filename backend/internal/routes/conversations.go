package routes

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"
	"strconv"
)

func handlerGetPrivate(w http.ResponseWriter, r *http.Request) {
	privateIdStr := r.PathValue("private_id")
	privateId, err := strconv.ParseInt(privateIdStr, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid private_id", nil)
		return
	}

	private, err := models.GetPrivateById(privateId)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while fetching private conversation", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Private conversation fetched successfully", private)

}