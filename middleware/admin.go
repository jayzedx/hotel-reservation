package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*repo.User)
	if !ok {
		logs.Info("user id is invalid")
		return unAuthorized
	}
	if !user.IsAdmin {
		logs.Info(fmt.Sprintf("%s : isn't in admin role", user.Email))
		return unAuthorized
	}
	return c.Next()
}
