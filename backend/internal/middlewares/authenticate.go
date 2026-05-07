package middlewares

import (
	"backend/internal/utils"
	"context"
	"net/http"
	"strings"
)

const (
	CtxUserID string = "userId"
	CtxName   string = "name"
	CtxPlatform string = "X-Platform"
	CtxAuthorization string = "Authorization"
	PlatformWeb string = "web"
	PlatformMobile string = "mobile"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.TrimSpace(r.Header.Get(CtxAuthorization))
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			utils.JSON(w, http.StatusUnauthorized, false, "Missing or invalid Authorization header", nil)
			return
		}

		platform := strings.TrimSpace(r.Header.Get(CtxPlatform))
		if platform != PlatformWeb && platform != PlatformMobile {
			utils.JSON(w, http.StatusBadRequest, false, "Missing or invalid X-Platform header", nil)
			return
		}

		accessToken := strings.TrimSpace(authHeader[7:]) // Remove "Bearer " prefix
		userId, name, tokenPlatform, err := utils.VerifyJWT(accessToken)
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, false, "Invalid or expired token: "+err.Error(), nil)
			return
		}
		if tokenPlatform != platform {
			utils.JSON(w, http.StatusUnauthorized, false, "Token platform does not match X-Platform header", nil)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, userId)
		ctx = context.WithValue(ctx, CtxName, name)
		ctx = context.WithValue(ctx, CtxPlatform, tokenPlatform)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}