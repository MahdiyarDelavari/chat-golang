package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func InitJWT(key string) {
	jwtKey = []byte(key)
}

type customClaims struct {
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	Platform string `json:"X-Platform"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId int64, name string, platform string) (string, error) {
	exp := time.Now().Add(2 * time.Hour)

	if platform != "web" && platform != "mobile" {
		return "", errors.New("invalid platform, must be 'web' or 'mobile'")
	}
	claims := customClaims{	
		UserID:   userId,
		Name:     name,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject: fmt.Sprint(userId),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenStr string) (int64,string,string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr, 
		&customClaims{},
		func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			return 0,"","", fmt.Errorf("token parse error: %v",err)
		}
		if !token.Valid {
			return 0,"","", errors.New("invalid token")
		}	
		claims, ok := token.Claims.(*customClaims)
		if !ok {
			return 0,"","", errors.New("invalid token claims")
		}
		if claims.UserID == 0 || claims.Name == "" || claims.Platform == "" {
			return 0,"","", errors.New("missing required claims")
		}

		return claims.UserID,claims.Name,claims.Platform, nil
}