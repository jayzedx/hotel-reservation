package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/db"
	"github.com/jayzedx/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *userHandler {
	return &userHandler{
		store: store,
	}
}

func (h *userHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *userHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.store.User.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *userHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.store.User.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)

}

func (h *userHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	if err := h.store.User.DeleteUser(c.Context(), oid); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(map[string]string{"deleted": userId})
}

func (h *userHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		// values bson.M
		params types.UpdateUserParams
		userId = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	// if err := c.BodyParser(&values); err != nil {
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.store.User.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userId})
}
