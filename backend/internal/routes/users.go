package routes

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"
	"strconv"
)

func handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")

	targetId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid user ID", nil)
		return
	}

	existingUser, err := models.GetUserByID(targetId)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, false, "User not found", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "User found", existingUser)

}