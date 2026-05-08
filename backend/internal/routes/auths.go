package routes

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
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

func handlerEmailLogin(w http.ResponseWriter, r *http.Request){
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(string(middlewares.CtxPlatform))))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Email and password are required", nil)
		return
	}

	existingUser, err := models.GetUserByEmail(req.Email)
	if err != nil || existingUser == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid email or password", nil)
		return
	}

	err = utils.CheckPasswordHash(req.Password, existingUser.Password) 
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid email or password", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(existingUser.ID, existingUser.Name, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while generating token", nil)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while generating refresh token", nil)
		return
	}

	err = models.UpdateUserRefreshToken(existingUser.ID, platform, refreshToken)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while saving refresh token", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Login successful", map[string]any{
		"user": existingUser,
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.CtxUserID).(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}
	platform, ok := r.Context().Value(middlewares.CtxPlatform).(string)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}
	err := models.DeleteUserRefreshToken(userID, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while logging out", nil)
		return
	}
	utils.JSON(w, http.StatusOK, true, "Logout successful", nil)

}

func handlerRefreshSession(w http.ResponseWriter, r *http.Request) {
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(string(middlewares.CtxPlatform))))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Refresh token is required", nil)
		return
	}

	existingUser, err := models.GetUserByRefreshToken(req.RefreshToken, platform)
	if err != nil || existingUser == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid credintials", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(existingUser.ID, existingUser.Name, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while generating token", nil)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while generating refresh token", nil)
		return
	}

	err = models.UpdateUserRefreshToken(existingUser.ID, platform, refreshToken)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Error occurred while saving refresh token", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Login successful", map[string]any{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func handlerGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(string(middlewares.CtxPlatform))))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Refresh token is required", nil)
		return
	}

	existingUser, err := models.GetUserByRefreshToken(req.RefreshToken, platform)
	if err != nil || existingUser == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid credintials", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "User fetched successfully", map[string]any{
		"user": existingUser,
	})
}