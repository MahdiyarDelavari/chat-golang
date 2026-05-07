package routes

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"net/http"
)

func handlerEmailRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Password string `json:"password"`
	}

	err :=json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Name, email and password are required", nil)
		return
	}

	existingUser , err := models.GetUserByEmail(req.Email)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while fetching user", nil)
		return
	}
	if existingUser != nil {
		utils.JSON(w, http.StatusConflict, false, "User with this email already exists", nil)
		return
	}

	hashedPassword , err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while hashing password", nil)
		return
	}

	user, err := models.CreateUserByEmail(req.Name, req.Email, hashedPassword)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while creating user", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "User created successfully", user)
}