package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/jayzedx/hotel-reservation/api/response"
	"github.com/jayzedx/hotel-reservation/db"
	"github.com/jayzedx/hotel-reservation/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct
type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:email`
	Password string `json:password`
}

type AuthResponse struct {
	User  *types.User `json:user`
	Token string      `json:token`
}

// Factory
func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

// Implement
func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return nil
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.InvalidCredentials(c, nil)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		// data := make(map[string]interface{})
		// users, ok := data["users"].([]*types.User)
		// if !ok {
		// 	users = make([]*types.User, 0)
		// }
		// data["users"] = append(users, &types.User{FirstName: "jay"})
		// return response.InvalidCredentials(c, &data)
		return response.InvalidCredentials(c, nil)
	}

	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}
	//claims["id"] = user.Id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := viper.GetString("app.jwt_secret")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenStr
}
