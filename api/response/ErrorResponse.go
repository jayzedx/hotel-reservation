package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func InvalidCredentials(c *fiber.Ctx, data *map[string]interface{}) error {
	return c.Status(http.StatusBadRequest).JSON(HeadResponse{
		StatusCode: http.StatusBadRequest,
		Status:     "error",
		Msg:        "invalid credentials",
		Data:       data,
	})
}

func UnAuthorized(c *fiber.Ctx, data *map[string]interface{}) error {
	return c.Status(http.StatusUnauthorized).JSON(HeadResponse{
		StatusCode: http.StatusUnauthorized,
		Status:     "error",
		Msg:        "unauthorized",
		Data:       data,
	})
}
