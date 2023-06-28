package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/db"
	"github.com/jayzedx/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct
type hotelHandler struct {
	store *db.Store
}

// function
func NewHotelHandler(store *db.Store) *hotelHandler {
	return &hotelHandler{
		store: store,
	}
}

// implement
type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *hotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	// fmt.Println(qparams.Rooms, qparams.Rating)

	query := bson.M{}
	if c.Query("rooms") != "" {
		query["rooms"] = bson.M{"$exists": qparams.Rooms}
	}
	if c.Query("rating") != "" {
		query["rating"] = qparams.Rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), query)

	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *hotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	hotels, err := h.store.Hotel.GetHotelById(c.Context(), oid)

	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *hotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelId": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	if len(rooms) == 0 {
		// make([]*types.Room, 0)
		return c.JSON([]*types.Room{})
	}
	return c.JSON(rooms)
}

func (h *hotelHandler) HandlePostHotel(c *fiber.Ctx) error {

	return nil
}

func (h *hotelHandler) HandlePutHotel(c *fiber.Ctx) error {

	return nil
}

func (h *hotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {

	return nil
}
