package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/spf13/viper"
)

type AuthService interface {
	Authenticate(params CreateAuthParams) (*AuthResponse, error)
}

type CreateAuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func MapAuthResponse(auth *repo.Auth, token string) *AuthResponse {
	return &AuthResponse{
		Email: auth.Email,
		Token: token,
	}
}

func CreateAuthFromUser(user *repo.User) *repo.Auth {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	return &repo.Auth{
		UserId:   user.Id,
		Email:    user.Email,
		Expires:  expires,
		CreateAt: now.Unix(),
	}
}

func createTokenFromAuth(auth *repo.Auth) (string, error) {
	claims := jwt.MapClaims{
		"id":      auth.Id,
		"user_id": auth.UserId,
		"email":   auth.Email,
		"expires": auth.Expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := viper.GetString("app.jwt_secret")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		logs.Error("Failed to sign token with secret")
		return "", err
	}
	// auth.Token = tokenStr
	return tokenStr, nil
}
