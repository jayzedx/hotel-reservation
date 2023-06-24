package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Jay",
		LastName:  "Layman",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("James")
}
