package models

import (
	"backend/internal/db"
	"backend/internal/middlewares"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"-"`
	RefreshTokenWeb      *string `json:"-"`
	RefreshTokenWebAt    *time.Time `json:"-"`
	RefreshTokenMobile   *string `json:"-"`
	RefreshTokenMobileAt *time.Time `json:"-"`
	CreatedAt            time.Time `json:"created_at"`
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	row := db.DB.QueryRow("SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at FROM users WHERE email = ?", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RefreshTokenWeb, &user.RefreshTokenWebAt, &user.RefreshTokenMobile, &user.RefreshTokenMobileAt, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUserByEmail(name, email, hashedPassword string) (*User, error) {
	result, err := db.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	return &User{
		ID: id,
		Name: name,
		Email: email,
		CreatedAt: createdAt,
	}, nil
}
func UpdateUserRefreshToken(userID int64, platform string, refreshToken string) error {
	now := time.Now()
	switch platform {
	case middlewares.PlatformWeb:
		_, err := db.DB.Exec("UPDATE users SET refresh_token_web = ?, refresh_token_web_at = ? WHERE id = ?", refreshToken, now, userID)
		return err
	case middlewares.PlatformMobile:
		_, err := db.DB.Exec("UPDATE users SET refresh_token_mobile = ?, refresh_token_mobile_at = ? WHERE id = ?", refreshToken, now, userID)
		return err
	default:
		return errors.New("invalid platform")
	}
}