package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/spf13/viper"
)

var unAuthorized = errs.AppError{
	Code:    http.StatusUnauthorized,
	Message: "Unauthorized",
}

func JWTAuthentication(c *fiber.Ctx) error {

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		logs.Info("Wrong request Header")
		return unAuthorized
	}

	claims, err := validateToken(token)

	if err != nil {
		logs.Info("Failed to validation JWT token : " + err.Error())
		return unAuthorized
	}

	// check token expired
	// fmt.Println("claims : ", claims)
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)
	if time.Now().Unix() > expires {
		logs.Info("Token Expired")
		return unAuthorized
	}
	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errStr := fmt.Sprint("Invalid siging method", token.Header["alg"]) //alg - allow origin
			logs.Info(errStr)
			return nil, unAuthorized
		}
		secret := viper.GetString("app.jwt_secret")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		logs.Info("Invalid token")
		return nil, unAuthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, unAuthorized
	}

	return claims, nil
}
