package routes

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
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

func handlerJoinPrivate(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.CtxUserID).(int64)

	var req struct {
		RecieverID int64 `json:"receiver_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}
	if req.RecieverID == 0 {
		utils.JSON(w, http.StatusBadRequest, false, "Receiver ID is required", nil)
		return
	}

	private, err := models.GetPrivateByUsers(userId, req.RecieverID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while fetching private conversation", nil)
		return
	}
	if private == nil {
		private, err = models.CreatePrivate(userId, req.RecieverID)
		if err != nil {
			utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while creating private conversation", nil)
			return
		}
		utils.JSON(w, http.StatusOK, true, "Joined private conversation successfully", private)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Joined private conversation successfully", private)
}

func handlerGetConversations(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.CtxUserID).(int64)

	privates, err := models.GetPrivatesForUser(userId)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while fetching conversations", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Conversations fetched successfully", privates)
}